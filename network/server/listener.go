package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Yeah114/FunShuttler/core/minecraft/protocol"
	fpacket "github.com/Yeah114/FunShuttler/core/minecraft/protocol/packet"
	"github.com/Yeah114/FunShuttler/game_control/game_interface"
	"github.com/Yeah114/FunShuttler/game_control/resources_control"
	"github.com/Yeah114/FunShuttler/uqholder"
	uqdefines "github.com/Yeah114/FunShuttler/uqholder/defines"
	"github.com/Yeah114/tempest-core/network/app"
	listenerpb "github.com/Yeah114/tempest-core/network_api/listener"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type packetEvent struct {
	packet fpacket.Packet
	err    error
}

type bytesEvent struct {
	id      uint32
	payload []byte
	err     error
}

// ListenerService streams packet and chat events to clients.
type ListenerService struct {
	listenerpb.UnimplementedListenerServiceServer
	state *app.FatalderState

	packetEvents chan packetEvent
	bytesEvents  chan bytesEvent

	packetStreamActive atomic.Bool
	bytesStreamActive  atomic.Bool

	mu                   sync.Mutex
	typedPacketListeners map[uint32]string
	typedBytesListeners  map[uint32]string
}

// NewListenerService constructs a listener service.
func NewListenerService(state *app.FatalderState) *ListenerService {
	svc := &ListenerService{
		state:                state,
		packetEvents:         make(chan packetEvent, 1024),
		bytesEvents:          make(chan bytesEvent, 1024),
		typedPacketListeners: make(map[uint32]string),
		typedBytesListeners:  make(map[uint32]string),
	}
	go svc.monitorDisconnects()
	return svc
}

func (s *ListenerService) monitorDisconnects() {
	for {
		ch, cancel := s.state.DisconnectEvents(1)
		_, ok := <-ch
		cancel()
		if !ok {
			continue
		}
		s.mu.Lock()
		s.typedPacketListeners = make(map[uint32]string)
		s.typedBytesListeners = make(map[uint32]string)
		s.mu.Unlock()
	}
}

func (s *ListenerService) ListenFateArk(req *listenerpb.ListenFateArkRequest, stream listenerpb.ListenerService_ListenFateArkServer) error {
	messages, cancel := s.state.Messages(64)
	defer cancel()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, ok := <-messages:
			if !ok {
				return nil
			}
			if err := stream.Send(&listenerpb.Output{
				MsgType: msg.Type,
				Msg:     msg.Message,
				ErrMsg:  msg.Error,
			}); err != nil {
				return err
			}
		}
	}
}

func (s *ListenerService) ListenPackets(req *listenerpb.ListenPacketsRequest, stream listenerpb.ListenerService_ListenPacketsServer) error {
	if !s.packetStreamActive.CompareAndSwap(false, true) {
		return status.Error(codes.ResourceExhausted, "packet stream already active")
	}
	defer s.packetStreamActive.Store(false)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt := <-s.packetEvents:
			if evt.err != nil {
				return evt.err
			}
			if evt.packet == nil {
				continue
			}
			payload, err := json.Marshal(evt.packet)
			if err != nil {
				return err
			}
			if err := stream.Send(&listenerpb.Packet{
				Id:      evt.packet.ID(),
				Payload: string(payload),
			}); err != nil {
				return err
			}
		}
	}
}

func (s *ListenerService) ListenBytesPackets(req *listenerpb.ListenBytesPacketsRequest, stream listenerpb.ListenerService_ListenBytesPacketsServer) error {
	if !s.bytesStreamActive.CompareAndSwap(false, true) {
		return status.Error(codes.ResourceExhausted, "bytes packet stream already active")
	}
	defer s.bytesStreamActive.Store(false)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt := <-s.bytesEvents:
			if evt.err != nil {
				return evt.err
			}
			if err := stream.Send(&listenerpb.BytesPacket{
				Id:      evt.id,
				Payload: evt.payload,
			}); err != nil {
				return err
			}
		}
	}
}

func (s *ListenerService) ListenTypedPacket(ctx context.Context, req *listenerpb.ListenTypedPacketRequest) (*responsepb.GeneralResponse, error) {
	packetID := req.GetPacketId()
	if packetID == 0 {
		return nil, status.Error(codes.InvalidArgument, "packet_id required")
	}

	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		pl := iface.PacketListener()
		if pl == nil {
			return errors.New("packet listener unavailable")
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		if _, exists := s.typedPacketListeners[packetID]; exists {
			return nil
		}
		listenerID, err := pl.ListenPacket([]uint32{packetID}, func(pk fpacket.Packet, connErr error) {
			s.pushPacketEvent(packetEvent{packet: pk, err: connErr})
		})
		if err != nil {
			return err
		}
		s.typedPacketListeners[packetID] = listenerID
		return nil
	})
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(""), nil
}

func (s *ListenerService) ListenTypedBytesPacket(ctx context.Context, req *listenerpb.ListenTypedBytesPacketRequest) (*responsepb.GeneralResponse, error) {
	packetID := req.GetPacketId()
	if packetID == 0 {
		return nil, status.Error(codes.InvalidArgument, "packet_id required")
	}

	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		pl := iface.PacketListener()
		if pl == nil {
			return errors.New("packet listener unavailable")
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		if _, exists := s.typedBytesListeners[packetID]; exists {
			return nil
		}
		listenerID, err := pl.ListenPacket([]uint32{packetID}, func(pk fpacket.Packet, connErr error) {
			if pk == nil {
				s.pushBytesEvent(bytesEvent{err: connErr})
				return
			}
			buf := bytes.NewBuffer(nil)
			writer := protocol.NewWriter(buf, 0)
			func() {
				defer func() {
					if recoverErr := recover(); recoverErr != nil {
						s.pushBytesEvent(bytesEvent{err: fmt.Errorf("marshal packet %d failed: %v", packetID, recoverErr)})
					}
				}()
				pk.Marshal(writer)
				s.pushBytesEvent(bytesEvent{payload: buf.Bytes(), id: pk.ID(), err: connErr})
			}()
		})
		if err != nil {
			return err
		}
		s.typedBytesListeners[packetID] = listenerID
		return nil
	})
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(""), nil
}

func (s *ListenerService) ListenPlayerChange(req *listenerpb.ListenPlayerChangeRequest, stream listenerpb.ListenerService_ListenPlayerChangeServer) error {
	type playerEvent struct {
		action string
		uuid   string
		err    error
	}

	events := make(chan playerEvent, 256)
	done := make(chan struct{})

	var listener *resources_control.PacketListener
	var listenerID string

	if err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		listener = iface.PacketListener()
		if listener == nil {
			return errors.New("packet listener unavailable")
		}
		var err error
		listenerID, err = listener.ListenPacket([]uint32{fpacket.IDPlayerList}, func(pk fpacket.Packet, connErr error) {
			select {
			case <-done:
				return
			default:
			}
			if connErr != nil {
				select {
				case events <- playerEvent{err: connErr}:
				default:
				}
				return
			}
			playerList, ok := pk.(*fpacket.PlayerList)
			if !ok {
				return
			}
			action := ""
			switch playerList.ActionType {
			case fpacket.PlayerListActionAdd:
				action = "online"
			case fpacket.PlayerListActionRemove:
				action = "offline"
			default:
				return
			}
			for _, entry := range playerList.Entries {
				select {
				case events <- playerEvent{action: action, uuid: entry.UUID.String()}:
				default:
				}
			}
		})
		return err
	}); err != nil {
		return err
	}
	defer func() {
		close(done)
		if listener != nil && listenerID != "" {
			listener.DestroyListener(listenerID)
		}
	}()

	// Emit existing players as "exist".
	registry := s.state.Players()
	players, err := s.state.SnapshotPlayers()
	if err == nil {
		for _, player := range players {
			if player == nil {
				continue
			}
			if uuidStr, ok := player.GetUUIDString(); ok {
				registry.Rebind(uuidStr, player)
			}
			if err := stream.Send(&listenerpb.PlayerAction{Action: "exist"}); err != nil {
				return err
			}
		}
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt := <-events:
			if evt.err != nil {
				return evt.err
			}
			if evt.action == "" {
				continue
			}
			if evt.uuid != "" {
				if evt.action == "offline" {
					registry.Delete(evt.uuid)
				} else {
					if player, lookupErr := s.lookupPlayer(evt.uuid); lookupErr == nil {
						registry.Rebind(evt.uuid, player)
					}
				}
			}
			if err := stream.Send(&listenerpb.PlayerAction{Action: evt.action}); err != nil {
				return err
			}
		}
	}
}

func (s *ListenerService) ListenChat(req *listenerpb.ListenChatRequest, stream listenerpb.ListenerService_ListenChatServer) error {
	return s.streamTextPackets("", stream)
}

func (s *ListenerService) ListenCommandBlock(req *listenerpb.ListenCommandBlockRequest, stream listenerpb.ListenerService_ListenCommandBlockServer) error {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return status.Error(codes.InvalidArgument, "name required")
	}
	return s.streamTextPackets(name, stream)
}

func (s *ListenerService) streamTextPackets(filter string, stream listenerpb.ListenerService_ListenChatServer) error {
	type chatEvent struct {
		payload []byte
		err     error
	}

	events := make(chan chatEvent, 256)
	done := make(chan struct{})

	var listener *resources_control.PacketListener
	var listenerID string

	if err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		listener = iface.PacketListener()
		if listener == nil {
			return errors.New("packet listener unavailable")
		}
		var err error
		listenerID, err = listener.ListenPacket([]uint32{fpacket.IDText}, func(pk fpacket.Packet, connErr error) {
			select {
			case <-done:
				return
			default:
			}
			if connErr != nil {
				select {
				case events <- chatEvent{err: connErr}:
				default:
				}
				return
			}
			text, ok := pk.(*fpacket.Text)
			if !ok {
				return
			}
			if filter != "" && text.SourceName != filter {
				return
			}
			payload, err := buildChatPayload(text)
			select {
			case events <- chatEvent{payload: payload, err: err}:
			default:
			}
		})
		return err
	}); err != nil {
		return err
	}
	defer func() {
		close(done)
		if listener != nil && listenerID != "" {
			listener.DestroyListener(listenerID)
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt := <-events:
			if evt.err != nil {
				return evt.err
			}
			if len(evt.payload) == 0 {
				continue
			}
			if err := stream.Send(&listenerpb.Chat{Payload: string(evt.payload)}); err != nil {
				return err
			}
		}
	}
}

func (s *ListenerService) pushPacketEvent(evt packetEvent) {
	select {
	case s.packetEvents <- evt:
	default:
	}
}

func (s *ListenerService) pushBytesEvent(evt bytesEvent) {
	select {
	case s.bytesEvents <- evt:
	default:
	}
}

func (s *ListenerService) lookupPlayer(uuidStr string) (uqdefines.PlayerUQReader, error) {
	return fetchPlayerByUUID(s.state, uuidStr)
}

type chatPayload struct {
	Name          string   `json:"name"`
	Msg           []string `json:"msg"`
	Type          byte     `json:"type"`
	RawMsg        string   `json:"raw_msg"`
	RawName       string   `json:"raw_name"`
	RawParameters []string `json:"raw_parameters"`
}

func buildChatPayload(text *fpacket.Text) ([]byte, error) {
	payload := chatPayload{
		Name:          uqholder.ToPlainName(text.SourceName),
		Msg:           splitWords(text.Message),
		Type:          text.TextType,
		RawMsg:        text.Message,
		RawName:       text.SourceName,
		RawParameters: append([]string(nil), text.Parameters...),
	}
	return json.Marshal(payload)
}

func splitWords(s string) []string {
	fields := strings.Fields(s)
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			result = append(result, field)
		}
	}
	return result
}
