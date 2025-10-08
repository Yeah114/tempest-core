package client

import (
	"context"
	"encoding/json"
	"strings"

	utilspb "github.com/Yeah114/tempest-core/network_api/utils"
	"google.golang.org/grpc"
)

type UtilsClient struct {
	rpc         utilspb.UtilsServiceClient
	callOptions []grpc.CallOption
}

func newUtilsClient(rpc utilspb.UtilsServiceClient, callOptions []grpc.CallOption) *UtilsClient {
	return &UtilsClient{
		rpc:         rpc,
		callOptions: append([]grpc.CallOption(nil), callOptions...),
	}
}

func (c *UtilsClient) ready() error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("utils")
	}
	return nil
}

func (c *UtilsClient) callOpts(opts []grpc.CallOption) []grpc.CallOption {
	return mergeCallOptions(c.callOptions, opts)
}

func (c *UtilsClient) SendPacket(ctx context.Context, packetID int32, jsonPayload string, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &utilspb.SendPacketRequest{
		PacketId: packetID,
		JsonStr:  strings.TrimSpace(jsonPayload),
	}
	resp, err := c.rpc.SendPacket(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *UtilsClient) GetPacketNameIDMapping(ctx context.Context, opts ...grpc.CallOption) (map[string]uint32, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	resp, err := c.rpc.GetPacketNameIDMapping(ctx, &utilspb.GetPacketNameIDMappingRequest{}, c.callOpts(opts)...)
	if err != nil {
		return nil, err
	}
	payload, err := generalPayload(resp)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload) == "" {
		return map[string]uint32{}, nil
	}
	out := make(map[string]uint32)
	if err := json.Unmarshal([]byte(payload), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *UtilsClient) GetClientMaintainedBotBasicInfo(ctx context.Context, opts ...grpc.CallOption) (map[string]any, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	resp, err := c.rpc.GetClientMaintainedBotBasicInfo(ctx, &utilspb.GetClientMaintainedBotBasicInfoRequest{}, c.callOpts(opts)...)
	if err != nil {
		return nil, err
	}
	payload, err := generalPayload(resp)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload) == "" {
		return map[string]any{}, nil
	}
	var info map[string]any
	if err := json.Unmarshal([]byte(payload), &info); err != nil {
		return nil, err
	}
	return info, nil
}

func (c *UtilsClient) GetClientMaintainedExtendInfo(ctx context.Context, opts ...grpc.CallOption) (map[string]any, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	resp, err := c.rpc.GetClientMaintainedExtendInfo(ctx, &utilspb.GetClientMaintainedExtendInfoRequest{}, c.callOpts(opts)...)
	if err != nil {
		return nil, err
	}
	payload, err := generalPayload(resp)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload) == "" {
		return map[string]any{}, nil
	}
	var info map[string]any
	if err := json.Unmarshal([]byte(payload), &info); err != nil {
		return nil, err
	}
	return info, nil
}
