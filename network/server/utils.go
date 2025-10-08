package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Yeah114/FunShuttler/game_control/resources_control"
	"github.com/Yeah114/tempest-core/network/app"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	utilspb "github.com/Yeah114/tempest-core/network_api/utils"
)

// UtilsService bridges misc helper endpoints.
type UtilsService struct {
	utilspb.UnimplementedUtilsServiceServer
	state *app.FatalderState
}

// NewUtilsService constructs a utils service.
func NewUtilsService(state *app.FatalderState) *UtilsService {
	return &UtilsService{state: state}
}

func (s *UtilsService) SendPacket(ctx context.Context, req *utilspb.SendPacketRequest) (*responsepb.GeneralResponse, error) {
	pool, err := s.state.PacketPool()
	if err != nil {
		return generalFailure(err), nil
	}
	ctor, ok := pool[uint32(req.GetPacketId())]
	if !ok || ctor == nil {
		return generalFailure(errors.New("packet id not supported")), nil
	}
	packet := ctor()
	if err := json.Unmarshal([]byte(req.GetJsonStr()), packet); err != nil {
		return generalFailure(err), nil
	}
	err = s.state.WithResources(func(res *resources_control.Resources) error {
		return res.WritePacket(packet)
	})
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(""), nil
}

func (s *UtilsService) GetPacketNameIDMapping(ctx context.Context, req *utilspb.GetPacketNameIDMappingRequest) (*responsepb.GeneralResponse, error) {
	mapping, err := s.state.PacketNameID()
	if err != nil {
		return generalFailure(err), nil
	}
	data, err := json.Marshal(mapping)
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(string(data)), nil
}

func (s *UtilsService) GetClientMaintainedBotBasicInfo(ctx context.Context, req *utilspb.GetClientMaintainedBotBasicInfoRequest) (*responsepb.GeneralResponse, error) {
	info := make(map[string]any)
	err := s.state.WithResources(func(res *resources_control.Resources) error {
		holder := res.UQHolder()
		if holder == nil {
			return errors.New("uqholder unavailable")
		}
		micro := holder.Micro()
		if micro == nil {
			return errors.New("micro uqholder unavailable")
		}
		basic := micro.GetBotBasicInfo()
		info["BotName"] = basic.GetBotName()
		info["BotRuntimeID"] = basic.GetBotRuntimeID()
		info["BotUniqueID"] = basic.GetBotUniqueID()
		info["BotIdentity"] = basic.GetBotIdentity()
		info["BotUUIDStr"] = basic.GetBotUUIDStr()
		return nil
	})
	if err != nil {
		return generalFailure(err), nil
	}
	data, err := json.Marshal(info)
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(string(data)), nil
}

func (s *UtilsService) GetClientMaintainedExtendInfo(ctx context.Context, req *utilspb.GetClientMaintainedExtendInfoRequest) (*responsepb.GeneralResponse, error) {
	info := make(map[string]any)
	err := s.state.WithResources(func(res *resources_control.Resources) error {
		holder := res.UQHolder()
		if holder == nil {
			return errors.New("uqholder unavailable")
		}
		micro := holder.Micro()
		if micro == nil {
			return errors.New("micro uqholder unavailable")
		}
		extend := micro.GetExtendInfo()
		if extend == nil {
			return errors.New("extend info unavailable")
		}
		if value, ok := extend.GetCompressThreshold(); ok {
			info["CompressThreshold"] = value
		}
		if value, ok := extend.GetWorldGameMode(); ok {
			info["WorldGameMode"] = value
		}
		if value, ok := extend.GetWorldDifficulty(); ok {
			info["WorldDifficulty"] = value
		}
		if value, ok := extend.GetTime(); ok {
			info["Time"] = value
		}
		if value, ok := extend.GetDayTime(); ok {
			info["DayTime"] = value
		}
		if value, ok := extend.GetDayTimePercent(); ok {
			info["TimePercent"] = value
		}
		if value, ok := extend.GetGameRules(); ok {
			info["GameRules"] = value
		}
		if value, ok := extend.GetCurrentTick(); ok {
			info["CurrentTick"] = value
		}
		if value, ok := extend.GetSyncRatio(); ok {
			info["SyncRatio"] = value
		}
		if value, ok := extend.GetBotDimension(); ok {
			info["BotDimension"] = value
		}
		if pos, syncTick := extend.GetBotPosition(); syncTick != 0 {
			info["BotPosition"] = []float32{pos[0], pos[1], pos[2]}
			info["BotPositionOutOfSyncTick"] = syncTick
		}
		info["ClientDimension"] = extend.GetClientDimension()
		info["ClientHotBarSlot"] = extend.GetClientHotBarSlot()
		info["ClientHoldingItem"] = extend.GetClientHoldingItem()
		return nil
	})
	if err != nil {
		return generalFailure(err), nil
	}
	data, err := json.Marshal(info)
	if err != nil {
		return generalFailure(err), nil
	}
	return generalSuccess(string(data)), nil
}
