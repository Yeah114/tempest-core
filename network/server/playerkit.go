package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Yeah114/FunShuttler/core/minecraft/protocol"
	fpacket "github.com/Yeah114/FunShuttler/core/minecraft/protocol/packet"
	"github.com/Yeah114/FunShuttler/game_control/game_interface"
	uqdefines "github.com/Yeah114/FunShuttler/uqholder/defines"
	"github.com/Yeah114/tempest-core/network/app"
	playerkitpb "github.com/Yeah114/tempest-core/network_api/playerkit"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PlayerKitService exposes player utilities over gRPC.
type PlayerKitService struct {
	playerkitpb.UnimplementedPlayerKitServiceServer
	state *app.FatalderState
}

// NewPlayerKitService constructs a new player kit service.
func NewPlayerKitService(state *app.FatalderState) *PlayerKitService {
	return &PlayerKitService{state: state}
}

func (s *PlayerKitService) GetAllOnlinePlayers(ctx context.Context, req *playerkitpb.GetAllOnlinePlayersRequest) (*responsepb.GeneralResponse, error) {
	players, err := s.state.SnapshotPlayers()
	if err != nil {
		return nil, toStatusError(err)
	}
	registry := s.state.Players()
	out := make([]string, 0, len(players))
	for _, player := range players {
		if player == nil {
			continue
		}
		uuidStr, ok := player.GetUUIDString()
		if !ok || uuidStr == "" {
			continue
		}
		registry.Rebind(uuidStr, player)
		out = append(out, uuidStr)
	}
	data, err := json.Marshal(out)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return generalSuccess(string(data)), nil
}

func (s *PlayerKitService) GetPlayerByName(ctx context.Context, req *playerkitpb.GetPlayerByNameRequest) (*responsepb.GeneralResponse, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "player name required")
	}

	player, err := fetchPlayerByName(s.state, name)
	if err != nil {
		return nil, toStatusError(err)
	}
	uuidStr, ok := player.GetUUIDString()
	if !ok || uuidStr == "" {
		return nil, status.Error(codes.Internal, "player uuid unavailable")
	}
	s.state.Players().Rebind(uuidStr, player)
	return generalSuccess(uuidStr), nil
}

func (s *PlayerKitService) GetPlayerByUUID(ctx context.Context, req *playerkitpb.GetPlayerByUUIDRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuid())
	if err != nil {
		return nil, toStatusError(err)
	}
	uuidStr, ok := player.GetUUIDString()
	if !ok || uuidStr == "" {
		return nil, status.Error(codes.Internal, "player uuid unavailable")
	}
	return generalSuccess(uuidStr), nil
}

func (s *PlayerKitService) ReleaseBindPlayer(ctx context.Context, req *playerkitpb.ReleaseBindPlayerRequest) (*responsepb.GeneralResponse, error) {
	s.state.Players().Delete(req.GetUuidStr())
	return generalSuccess(""), nil
}

func (s *PlayerKitService) GetPlayerName(ctx context.Context, req *playerkitpb.GetPlayerNameRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	name, ok := player.GetUsername()
	if !ok || name == "" {
		return nil, status.Error(codes.Internal, "player name unavailable")
	}
	return generalSuccess(name), nil
}

func (s *PlayerKitService) GetPlayerEntityUniqueID(ctx context.Context, req *playerkitpb.GetPlayerEntityUniqueIDRequest) (*responsepb.GeneralInt64Response, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	id, ok := player.GetEntityUniqueID()
	if !ok {
		return nil, status.Error(codes.Internal, "entity unique id unavailable")
	}
	return int64Success(id), nil
}

func (s *PlayerKitService) GetPlayerLoginTime(ctx context.Context, req *playerkitpb.GetPlayerLoginTimeRequest) (*responsepb.GeneralInt64Response, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	ts, ok := player.GetLoginTime()
	if !ok {
		return nil, status.Error(codes.Internal, "login time unavailable")
	}
	return int64Success(ts.Unix()), nil
}

func (s *PlayerKitService) GetPlayerPlatformChatID(ctx context.Context, req *playerkitpb.GetPlayerPlatformChatIDRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	value, ok := player.GetPlatformChatID()
	if !ok {
		return nil, status.Error(codes.Internal, "platform chat id unavailable")
	}
	return generalSuccess(value), nil
}

func (s *PlayerKitService) GetPlayerBuildPlatform(ctx context.Context, req *playerkitpb.GetPlayerBuildPlatformRequest) (*responsepb.GeneralInt32Response, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	value, ok := player.GetBuildPlatform()
	if !ok {
		return nil, status.Error(codes.Internal, "build platform unavailable")
	}
	return int32Success(value), nil
}

func (s *PlayerKitService) GetPlayerSkinID(ctx context.Context, req *playerkitpb.GetPlayerSkinIDRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	value, ok := player.GetSkinID()
	if !ok {
		return nil, status.Error(codes.Internal, "skin id unavailable")
	}
	return generalSuccess(value), nil
}

func (s *PlayerKitService) GetPlayerCanBuild(ctx context.Context, req *playerkitpb.GetPlayerCanBuildRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityBuild)
}

func (s *PlayerKitService) SetPlayerCanBuild(ctx context.Context, req *playerkitpb.SetPlayerCanBuildRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityBuild, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanDig(ctx context.Context, req *playerkitpb.GetPlayerCanDigRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityMine)
}

func (s *PlayerKitService) SetPlayerCanDig(ctx context.Context, req *playerkitpb.SetPlayerCanDigRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityMine, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanDoorsAndSwitches(ctx context.Context, req *playerkitpb.GetPlayerCanDoorsAndSwitchesRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityDoorsAndSwitches)
}

func (s *PlayerKitService) SetPlayerCanDoorsAndSwitches(ctx context.Context, req *playerkitpb.SetPlayerCanDoorsAndSwitchesRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityDoorsAndSwitches, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanOpenContainers(ctx context.Context, req *playerkitpb.GetPlayerCanOpenContainersRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityOpenContainers)
}

func (s *PlayerKitService) SetPlayerCanOpenContainers(ctx context.Context, req *playerkitpb.SetPlayerCanOpenContainersRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityOpenContainers, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanAttackPlayers(ctx context.Context, req *playerkitpb.GetPlayerCanAttackPlayersRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityAttackPlayers)
}

func (s *PlayerKitService) SetPlayerCanAttackPlayers(ctx context.Context, req *playerkitpb.SetPlayerCanAttackPlayersRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityAttackPlayers, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanAttackMobs(ctx context.Context, req *playerkitpb.GetPlayerCanAttackMobsRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityAttackMobs)
}

func (s *PlayerKitService) SetPlayerCanAttackMobs(ctx context.Context, req *playerkitpb.SetPlayerCanAttackMobsRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityAttackMobs, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanOperatorCommands(ctx context.Context, req *playerkitpb.GetPlayerCanOperatorCommandsRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityOperatorCommands)
}

func (s *PlayerKitService) SetPlayerCanOperatorCommands(ctx context.Context, req *playerkitpb.SetPlayerCanOperatorCommandsRequest) (*responsepb.GeneralResponse, error) {
	return s.setAbility(req.GetUuidStr(), protocol.AbilityOperatorCommands, req.GetAllow())
}

func (s *PlayerKitService) GetPlayerCanTeleport(ctx context.Context, req *playerkitpb.GetPlayerCanTeleportRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityTeleport)
}

func (s *PlayerKitService) SetPlayerCanTeleport(ctx context.Context, req *playerkitpb.SetPlayerCanTeleportRequest) (*responsepb.GeneralBoolResponse, error) {
	if _, err := s.setAbility(req.GetUuidStr(), protocol.AbilityTeleport, req.GetAllow()); err != nil {
		return nil, err
	}
	return boolSuccess(true), nil
}

func (s *PlayerKitService) GetPlayerStatusInvulnerable(ctx context.Context, req *playerkitpb.GetPlayerStatusInvulnerableRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityInvulnerable)
}

func (s *PlayerKitService) GetPlayerStatusFlying(ctx context.Context, req *playerkitpb.GetPlayerStatusFlyingRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityFlying)
}

func (s *PlayerKitService) GetPlayerStatusMayFly(ctx context.Context, req *playerkitpb.GetPlayerStatusMayFlyRequest) (*responsepb.GeneralBoolResponse, error) {
	return s.abilityBool(req.GetUuidStr(), protocol.AbilityMayFly)
}

func (s *PlayerKitService) GetPlayerDeviceID(ctx context.Context, req *playerkitpb.GetPlayerDeviceIDRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	value, ok := player.GetDeviceID()
	if !ok {
		return nil, status.Error(codes.Internal, "device id unavailable")
	}
	return generalSuccess(value), nil
}

func (s *PlayerKitService) GetPlayerEntityRuntimeID(ctx context.Context, req *playerkitpb.GetPlayerEntityRuntimeIDRequest) (*responsepb.GeneralUint64Response, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	value, ok := player.GetEntityRuntimeID()
	if !ok {
		return nil, status.Error(codes.Internal, "runtime id unavailable")
	}
	return uint64Success(value), nil
}

func (s *PlayerKitService) GetPlayerEntityMetadata(ctx context.Context, req *playerkitpb.GetPlayerEntityMetadataRequest) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	meta, ok := player.GetEntityMetadata()
	if !ok {
		return nil, status.Error(codes.Internal, "entity metadata unavailable")
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return generalSuccess(string(data)), nil
}

func (s *PlayerKitService) GetPlayerIsOP(ctx context.Context, req *playerkitpb.GetPlayerIsOPRequest) (*responsepb.GeneralBoolResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	if player == nil {
		return nil, status.Error(codes.Internal, "player not bound")
	}
	if perm, ok := player.GetCommandPermissions(); ok {
		if perm >= byte(fpacket.CommandPermissionLevelAdmin) {
			return boolSuccess(true), nil
		}
	}
	value, abErr := abilityEnabled(player, protocol.AbilityOperatorCommands)
	if abErr != nil {
		return nil, toStatusError(abErr)
	}
	return boolSuccess(value), nil
}

func (s *PlayerKitService) GetPlayerOnline(ctx context.Context, req *playerkitpb.GetPlayerOnlineRequest) (*responsepb.GeneralBoolResponse, error) {
	player, err := s.ensurePlayer(req.GetUuidStr())
	if err != nil {
		return nil, toStatusError(err)
	}
	return boolSuccess(player.StillOnline()), nil
}

func (s *PlayerKitService) SendPlayerChat(ctx context.Context, req *playerkitpb.SendPlayerChatRequest) (*responsepb.GeneralResponse, error) {
	return s.sendMessage(req.GetUuidStr(), strings.TrimSpace(req.GetMsg()), "tellraw", "")
}

func (s *PlayerKitService) SendPlayerRawChat(ctx context.Context, req *playerkitpb.SendPlayerRawChatRequest) (*responsepb.GeneralResponse, error) {
	return s.sendMessage(req.GetUuidStr(), strings.TrimSpace(req.GetMsg()), "tellraw", "")
}

func (s *PlayerKitService) SendPlayerTitle(ctx context.Context, req *playerkitpb.SendPlayerTitleRequest) (*responsepb.GeneralResponse, error) {
	uuid := req.GetUuidStr()
	title := strings.TrimSpace(req.GetTitle())
	subTitle := strings.TrimSpace(req.GetSubTitle())
	player, err := s.ensurePlayer(uuid)
	if err != nil {
		return nil, toStatusError(err)
	}
	name, ok := player.GetUsername()
	if !ok || name == "" {
		return nil, status.Error(codes.Internal, "player name unavailable")
	}

	err = s.withCommands(func(cmds *game_interface.Commands) error {
		target := quotedCommandTarget(name)
		if title != "" {
			payload, err := buildRawText(title)
			if err != nil {
				return err
			}
			if sendErr := cmds.SendWSCommand(fmt.Sprintf("titleraw %s title %s", target, payload)); sendErr != nil {
				return sendErr
			}
		}
		if subTitle != "" {
			payload, err := buildRawText(subTitle)
			if err != nil {
				return err
			}
			if sendErr := cmds.SendWSCommand(fmt.Sprintf("titleraw %s subtitle %s", target, payload)); sendErr != nil {
				return sendErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *PlayerKitService) SendPlayerActionBar(ctx context.Context, req *playerkitpb.SendPlayerActionBarRequest) (*responsepb.GeneralResponse, error) {
	return s.sendMessage(req.GetUuidStr(), strings.TrimSpace(req.GetActionBar()), "titleraw", "actionbar")
}

func (s *PlayerKitService) InterceptPlayerJustNextInput(ctx context.Context, req *playerkitpb.InterceptPlayerJustNextInputRequest) (*responsepb.GeneralResponse, error) {
	// Not yet implemented; acknowledge the call.
	return generalSuccess(""), nil
}

func (s *PlayerKitService) ensurePlayer(uuidStr string) (uqdefines.PlayerUQReader, error) {
	if uuidStr == "" {
		return nil, status.Error(codes.InvalidArgument, "player uuid required")
	}
	registry := s.state.Players()
	if player, ok := registry.Get(uuidStr); ok && player != nil {
		return player, nil
	}
	player, err := fetchPlayerByUUID(s.state, uuidStr)
	if err != nil {
		return nil, err
	}
	registry.Rebind(uuidStr, player)
	return player, nil
}

func (s *PlayerKitService) abilityBool(uuid string, ability uint32) (*responsepb.GeneralBoolResponse, error) {
	player, err := s.ensurePlayer(uuid)
	if err != nil {
		return nil, toStatusError(err)
	}
	value, abErr := abilityEnabled(player, ability)
	if abErr != nil {
		return nil, toStatusError(abErr)
	}
	return boolSuccess(value), nil
}

func (s *PlayerKitService) setAbility(uuid string, ability uint32, allow bool) (*responsepb.GeneralResponse, error) {
	player, err := s.ensurePlayer(uuid)
	if err != nil {
		return nil, toStatusError(err)
	}
	if err := updateAbility(s.state, player, ability, allow); err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *PlayerKitService) sendMessage(uuid, message, command, action string) (*responsepb.GeneralResponse, error) {
	if message == "" {
		return nil, status.Error(codes.InvalidArgument, "message required")
	}
	player, err := s.ensurePlayer(uuid)
	if err != nil {
		return nil, toStatusError(err)
	}
	name, ok := player.GetUsername()
	if !ok || name == "" {
		return nil, status.Error(codes.Internal, "player name unavailable")
	}

	err = s.withCommands(func(cmds *game_interface.Commands) error {
		payload, err := buildRawText(message)
		if err != nil {
			return err
		}
		target := quotedCommandTarget(name)
		cmd := command
		if action != "" {
			cmd = fmt.Sprintf("%s %s %s %s", command, target, action, payload)
		} else {
			cmd = fmt.Sprintf("%s %s %s", command, target, payload)
		}
		return cmds.SendWSCommand(cmd)
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *PlayerKitService) withCommands(fn func(*game_interface.Commands) error) error {
	return s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		cmds := iface.Commands()
		if cmds == nil {
			return status.Error(codes.FailedPrecondition, "commands interface unavailable")
		}
		return fn(cmds)
	})
}
