package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Yeah114/Fatalder/control"
	"github.com/Yeah114/Fatalder/defines"
	"github.com/Yeah114/FunShuttler/core/minecraft/protocol/packet"
	"github.com/Yeah114/FunShuttler/game_control/game_interface"
	"github.com/Yeah114/FunShuttler/game_control/resources_control"
	uqdefines "github.com/Yeah114/FunShuttler/uqholder/defines"
)

var (
	// ErrNotConnected is returned when an operation requires an active connection.
	ErrNotConnected = control.ErrNotConnected
	// ErrAlreadyConnected mirrors the control package sentinel.
	ErrAlreadyConnected = control.ErrAlreadyConnected
)

// ConnectOptions describes the parameters required to dial a rental server.
type ConnectOptions struct {
	AuthServerAddress string
	AuthUsername      string
	AuthPassword      string
	AuthToken         string
	ServerCode        string
	ServerPassword    string
}

// FatalderState tracks the shared connection state for the gRPC services.
type FatalderState struct {
	mu sync.RWMutex

	ctrl       *control.Control
	resources  *resources_control.Resources
	gameIface  *game_interface.GameInterface
	packetPool packet.Pool

	packetNameID map[string]uint32
	packetIDName map[uint32]string

	messageBus    *Broadcast[Message]
	disconnectBus *Broadcast[error]

	players *PlayerRegistry
}

// NewFatalderState creates a ready state container.
func NewFatalderState() *FatalderState {
	return &FatalderState{
		messageBus:    NewBroadcast[Message](),
		disconnectBus: NewBroadcast[error](),
		players:       NewPlayerRegistry(),
	}
}

// Connect establishes a new Control-driven session.
func (s *FatalderState) Connect(ctx context.Context, opts ConnectOptions) error {
	if opts.ServerCode == "" {
		return errors.New("server code required")
	}

	cfg := defines.ControlConfig{
		AuthServerAddress:    opts.AuthServerAddress,
		AuthServerToken:      opts.AuthToken,
		AuthUsername:         opts.AuthUsername,
		AuthPassword:         opts.AuthPassword,
		RentalServerCode:     opts.ServerCode,
		RentalServerPasscode: opts.ServerPassword,
	}

	ctrl := control.NewControl(cfg)
	if err := ctrl.EnterRentalServer(); err != nil {
		return err
	}

	resources := ctrl.Resources()
	gameIface := ctrl.GameInterface()
	if resources == nil || gameIface == nil {
		_ = ctrl.LeaveRentalServer()
		return errors.New("fatalder control returned nil resources")
	}

	packetPool := packet.ListAllPackets()
	nameID, idName := buildPacketMappings(packetPool)

	connCtx := ctrl.Client().Conn().Context()

	s.mu.Lock()
	if s.ctrl != nil {
		s.mu.Unlock()
		_ = ctrl.LeaveRentalServer()
		return ErrAlreadyConnected
	}

	// reset disconnect bus to ensure fresh subscriptions
	if s.disconnectBus != nil {
		s.disconnectBus.Close()
	}
	s.disconnectBus = NewBroadcast[error]()
	s.players = NewPlayerRegistry()

	s.ctrl = ctrl
	s.resources = resources
	s.gameIface = gameIface
	s.packetPool = packetPool
	s.packetNameID = nameID
	s.packetIDName = idName
	s.mu.Unlock()

	s.publishMessage(Message{
		Type:      "status",
		Message:   "connected",
		Timestamp: time.Now(),
	})

	go s.watchDisconnect(connCtx)

	return nil
}

// Disconnect tears down an active session.
func (s *FatalderState) Disconnect() error {
	s.mu.Lock()
	ctrl := s.ctrl
	if ctrl == nil {
		s.mu.Unlock()
		return ErrNotConnected
	}
	s.ctrl = nil
	s.resources = nil
	s.gameIface = nil
	s.packetPool = nil
	s.packetNameID = nil
	s.packetIDName = nil
	disconnectBus := s.disconnectBus
	s.disconnectBus = NewBroadcast[error]()
	s.players = NewPlayerRegistry()
	s.mu.Unlock()

	if disconnectBus != nil {
		disconnectBus.Publish(context.Canceled)
		disconnectBus.Close()
	}

	if err := ctrl.LeaveRentalServer(); err != nil {
		return err
	}

	s.publishMessage(Message{
		Type:      "status",
		Message:   "disconnected",
		Timestamp: time.Now(),
	})

	return nil
}

func (s *FatalderState) watchDisconnect(ctx context.Context) {
	<-ctx.Done()
	err := context.Cause(ctx)

	s.mu.Lock()
	ctrl := s.ctrl
	s.ctrl = nil
	s.resources = nil
	s.gameIface = nil
	s.packetPool = nil
	s.packetNameID = nil
	s.packetIDName = nil

	disconnectBus := s.disconnectBus
	s.disconnectBus = NewBroadcast[error]()
	s.players = NewPlayerRegistry()
	s.mu.Unlock()

	if ctrl != nil {
		_ = ctrl.LeaveRentalServer()
	}

	if disconnectBus != nil {
		disconnectBus.Publish(err)
		disconnectBus.Close()
	}

	s.publishMessage(Message{
		Type:      "disconnect",
		Message:   "connection closed",
		Error:     errString(err),
		Timestamp: time.Now(),
	})
}

// WithGameInterface executes fn while holding a read lock on the active game interface.
func (s *FatalderState) WithGameInterface(fn func(*game_interface.GameInterface) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.gameIface == nil {
		return ErrNotConnected
	}
	return fn(s.gameIface)
}

// WithResources executes fn with the resources guard.
func (s *FatalderState) WithResources(fn func(*resources_control.Resources) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.resources == nil {
		return ErrNotConnected
	}
	return fn(s.resources)
}

// PacketNameID returns the mapping from packet struct name to ID.
func (s *FatalderState) PacketNameID() (map[string]uint32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.packetNameID == nil {
		return nil, ErrNotConnected
	}
	out := make(map[string]uint32, len(s.packetNameID))
	for k, v := range s.packetNameID {
		out[k] = v
	}
	return out, nil
}

// PacketPool exposes the packet factory pool.
func (s *FatalderState) PacketPool() (packet.Pool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.packetPool == nil {
		return nil, ErrNotConnected
	}
	return s.packetPool, nil
}

// Players returns the current player registry. It is always non-nil.
func (s *FatalderState) Players() *PlayerRegistry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.players
}

// Messages channel for general updates.
func (s *FatalderState) Messages(buffer int) (<-chan Message, func()) {
	return s.messageBus.Subscribe(buffer)
}

// DisconnectEvents yields connection termination notifications.
func (s *FatalderState) DisconnectEvents(buffer int) (<-chan error, func()) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.disconnectBus.Subscribe(buffer)
}

// publishMessage is a helper to push to the broadcast.
func (s *FatalderState) publishMessage(msg Message) {
	if s.messageBus == nil {
		return
	}
	s.messageBus.Publish(msg)
}

func buildPacketMappings(pool packet.Pool) (map[string]uint32, map[uint32]string) {
	nameID := make(map[string]uint32, len(pool))
	idName := make(map[uint32]string, len(pool))
	for id, ctor := range pool {
		if ctor == nil {
			continue
		}
		pk := ctor()
		if pk == nil {
			continue
		}
		t := fmt.Sprintf("%T", pk)
		// Trim pointer and package prefix for readability.
		if len(t) > 0 {
			if t[0] == '*' {
				t = t[1:]
			}
			if idx := lastDot(t); idx >= 0 && idx < len(t)-1 {
				t = t[idx+1:]
			}
		}
		nameID[t] = id
		idName[id] = t
	}
	return nameID, idName
}

func lastDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// SnapshotPlayers copies all known players into the registry and returns them.
func (s *FatalderState) SnapshotPlayers() ([]uqdefines.PlayerUQReader, error) {
	var players []uqdefines.PlayerUQReader
	err := s.WithResources(func(res *resources_control.Resources) error {
		uq := res.UQHolder()
		if uq == nil {
			return errors.New("uqholder unavailable")
		}
		micro := uq.Micro()
		if micro == nil {
			return errors.New("micro uqholder unavailable")
		}
		players = micro.GetAllOnlinePlayers()
		registry := s.Players()
		for _, p := range players {
			if p == nil {
				continue
			}
			if uuidStr, ok := p.GetUUIDString(); ok {
				registry.Rebind(uuidStr, p)
			}
		}
		return nil
	})
	return players, err
}
