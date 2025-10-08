package client

import (
	"context"
	"encoding/json"
	"strings"

	playerkitpb "github.com/Yeah114/tempest-core/network_api/playerkit"
	"google.golang.org/grpc"
)

type PlayerKitClient struct {
	rpc         playerkitpb.PlayerKitServiceClient
	callOptions []grpc.CallOption
}

func newPlayerKitClient(rpc playerkitpb.PlayerKitServiceClient, callOptions []grpc.CallOption) *PlayerKitClient {
	return &PlayerKitClient{
		rpc:         rpc,
		callOptions: append([]grpc.CallOption(nil), callOptions...),
	}
}

func (c *PlayerKitClient) ready() error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("playerkit")
	}
	return nil
}

func (c *PlayerKitClient) callOpts(opts []grpc.CallOption) []grpc.CallOption {
	return mergeCallOptions(c.callOptions, opts)
}

func (c *PlayerKitClient) GetAllOnlinePlayers(ctx context.Context, opts ...grpc.CallOption) ([]string, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	resp, err := c.rpc.GetAllOnlinePlayers(ctx, &playerkitpb.GetAllOnlinePlayersRequest{}, c.callOpts(opts)...)
	if err != nil {
		return nil, err
	}
	payload, err := generalPayload(resp)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload) == "" {
		return nil, nil
	}
	var players []string
	if err := json.Unmarshal([]byte(payload), &players); err != nil {
		return nil, err
	}
	return players, nil
}

func (c *PlayerKitClient) GetPlayerByName(ctx context.Context, name string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	req := &playerkitpb.GetPlayerByNameRequest{Name: strings.TrimSpace(name)}
	resp, err := c.rpc.GetPlayerByName(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerByUUID(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	req := &playerkitpb.GetPlayerByUUIDRequest{Uuid: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerByUUID(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) ReleaseBindPlayer(ctx context.Context, uuid string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.ReleaseBindPlayerRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.ReleaseBindPlayer(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerName(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	req := &playerkitpb.GetPlayerNameRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerName(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerEntityUniqueID(ctx context.Context, uuid string, opts ...grpc.CallOption) (int64, error) {
	if err := c.ready(); err != nil {
		return 0, err
	}
	req := &playerkitpb.GetPlayerEntityUniqueIDRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerEntityUniqueID(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return 0, err
	}
	return int64Payload(resp)
}

func (c *PlayerKitClient) GetPlayerLoginTime(ctx context.Context, uuid string, opts ...grpc.CallOption) (int64, error) {
	if err := c.ready(); err != nil {
		return 0, err
	}
	req := &playerkitpb.GetPlayerLoginTimeRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerLoginTime(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return 0, err
	}
	return int64Payload(resp)
}

func (c *PlayerKitClient) GetPlayerPlatformChatID(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	req := &playerkitpb.GetPlayerPlatformChatIDRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerPlatformChatID(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerBuildPlatform(ctx context.Context, uuid string, opts ...grpc.CallOption) (int32, error) {
	if err := c.ready(); err != nil {
		return 0, err
	}
	req := &playerkitpb.GetPlayerBuildPlatformRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerBuildPlatform(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return 0, err
	}
	return int32Payload(resp)
}

func (c *PlayerKitClient) GetPlayerSkinID(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	req := &playerkitpb.GetPlayerSkinIDRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerSkinID(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerCanBuild(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	req := &playerkitpb.GetPlayerCanBuildRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerCanBuild(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanBuild(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanBuildRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanBuild(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanDig(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	req := &playerkitpb.GetPlayerCanDigRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerCanDig(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanDig(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanDigRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanDig(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanDoorsAndSwitches(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	req := &playerkitpb.GetPlayerCanDoorsAndSwitchesRequest{UuidStr: strings.TrimSpace(uuid)}
	resp, err := c.rpc.GetPlayerCanDoorsAndSwitches(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanDoorsAndSwitches(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanDoorsAndSwitchesRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanDoorsAndSwitches(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanOpenContainers(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerCanOpenContainers(ctx, &playerkitpb.GetPlayerCanOpenContainersRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanOpenContainers(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanOpenContainersRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanOpenContainers(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanAttackPlayers(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerCanAttackPlayers(ctx, &playerkitpb.GetPlayerCanAttackPlayersRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanAttackPlayers(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanAttackPlayersRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanAttackPlayers(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanAttackMobs(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerCanAttackMobs(ctx, &playerkitpb.GetPlayerCanAttackMobsRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanAttackMobs(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanAttackMobsRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanAttackMobs(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanOperatorCommands(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerCanOperatorCommands(ctx, &playerkitpb.GetPlayerCanOperatorCommandsRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanOperatorCommands(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SetPlayerCanOperatorCommandsRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanOperatorCommands(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) GetPlayerCanTeleport(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerCanTeleport(ctx, &playerkitpb.GetPlayerCanTeleportRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SetPlayerCanTeleport(ctx context.Context, uuid string, allow bool, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	req := &playerkitpb.SetPlayerCanTeleportRequest{UuidStr: strings.TrimSpace(uuid), Allow: allow}
	resp, err := c.rpc.SetPlayerCanTeleport(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) GetPlayerStatusInvulnerable(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerStatusInvulnerable(ctx, &playerkitpb.GetPlayerStatusInvulnerableRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) GetPlayerStatusFlying(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerStatusFlying(ctx, &playerkitpb.GetPlayerStatusFlyingRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) GetPlayerStatusMayFly(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerStatusMayFly(ctx, &playerkitpb.GetPlayerStatusMayFlyRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) GetPlayerDeviceID(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	resp, err := c.rpc.GetPlayerDeviceID(ctx, &playerkitpb.GetPlayerDeviceIDRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerEntityRuntimeID(ctx context.Context, uuid string, opts ...grpc.CallOption) (uint64, error) {
	if err := c.ready(); err != nil {
		return 0, err
	}
	resp, err := c.rpc.GetPlayerEntityRuntimeID(ctx, &playerkitpb.GetPlayerEntityRuntimeIDRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return 0, err
	}
	return uint64Payload(resp)
}

func (c *PlayerKitClient) GetPlayerEntityMetadata(ctx context.Context, uuid string, opts ...grpc.CallOption) (string, error) {
	if err := c.ready(); err != nil {
		return "", err
	}
	resp, err := c.rpc.GetPlayerEntityMetadata(ctx, &playerkitpb.GetPlayerEntityMetadataRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *PlayerKitClient) GetPlayerIsOP(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerIsOP(ctx, &playerkitpb.GetPlayerIsOPRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) GetPlayerOnline(ctx context.Context, uuid string, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.GetPlayerOnline(ctx, &playerkitpb.GetPlayerOnlineRequest{UuidStr: strings.TrimSpace(uuid)}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return boolPayload(resp)
}

func (c *PlayerKitClient) SendPlayerChat(ctx context.Context, uuid, message string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SendPlayerChatRequest{UuidStr: strings.TrimSpace(uuid), Msg: strings.TrimSpace(message)}
	resp, err := c.rpc.SendPlayerChat(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) SendPlayerRawChat(ctx context.Context, uuid, message string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SendPlayerRawChatRequest{UuidStr: strings.TrimSpace(uuid), Msg: strings.TrimSpace(message)}
	resp, err := c.rpc.SendPlayerRawChat(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) SendPlayerTitle(ctx context.Context, uuid, title, subTitle string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SendPlayerTitleRequest{
		UuidStr:  strings.TrimSpace(uuid),
		Title:    strings.TrimSpace(title),
		SubTitle: strings.TrimSpace(subTitle),
	}
	resp, err := c.rpc.SendPlayerTitle(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) SendPlayerActionBar(ctx context.Context, uuid, actionBar string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.SendPlayerActionBarRequest{UuidStr: strings.TrimSpace(uuid), ActionBar: strings.TrimSpace(actionBar)}
	resp, err := c.rpc.SendPlayerActionBar(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *PlayerKitClient) InterceptPlayerJustNextInput(ctx context.Context, uuid, retrieverID string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &playerkitpb.InterceptPlayerJustNextInputRequest{
		UuidStr:     strings.TrimSpace(uuid),
		RetrieverId: strings.TrimSpace(retrieverID),
	}
	resp, err := c.rpc.InterceptPlayerJustNextInput(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}
