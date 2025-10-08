package client

import (
	"context"

	listenerpb "github.com/Yeah114/tempest-core/network_api/listener"
	"google.golang.org/grpc"
)

type ListenerClient struct {
	rpc         listenerpb.ListenerServiceClient
	callOptions []grpc.CallOption
}

func newListenerClient(rpc listenerpb.ListenerServiceClient, callOptions []grpc.CallOption) *ListenerClient {
	return &ListenerClient{
		rpc:         rpc,
		callOptions: append([]grpc.CallOption(nil), callOptions...),
	}
}

func (c *ListenerClient) ready() error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("listener")
	}
	return nil
}

func (c *ListenerClient) callOpts(opts []grpc.CallOption) []grpc.CallOption {
	return mergeCallOptions(c.callOptions, opts)
}

func (c *ListenerClient) ListenTypedPacket(ctx context.Context, packetID uint32, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &listenerpb.ListenTypedPacketRequest{PacketId: packetID}
	resp, err := c.rpc.ListenTypedPacket(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *ListenerClient) ListenTypedBytesPacket(ctx context.Context, packetID uint32, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &listenerpb.ListenTypedBytesPacketRequest{PacketId: packetID}
	resp, err := c.rpc.ListenTypedBytesPacket(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *ListenerClient) ListenFateArk(ctx context.Context, req *listenerpb.ListenFateArkRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenFateArkClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenFateArkRequest{}
	}
	return c.rpc.ListenFateArk(ctx, req, c.callOpts(opts)...)
}

func (c *ListenerClient) ListenPackets(ctx context.Context, req *listenerpb.ListenPacketsRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenPacketsClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenPacketsRequest{}
	}
	return c.rpc.ListenPackets(ctx, req, c.callOpts(opts)...)
}

func (c *ListenerClient) ListenBytesPackets(ctx context.Context, req *listenerpb.ListenBytesPacketsRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenBytesPacketsClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenBytesPacketsRequest{}
	}
	return c.rpc.ListenBytesPackets(ctx, req, c.callOpts(opts)...)
}

func (c *ListenerClient) ListenPlayerChange(ctx context.Context, req *listenerpb.ListenPlayerChangeRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenPlayerChangeClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenPlayerChangeRequest{}
	}
	return c.rpc.ListenPlayerChange(ctx, req, c.callOpts(opts)...)
}

func (c *ListenerClient) ListenChat(ctx context.Context, req *listenerpb.ListenChatRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenChatClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenChatRequest{}
	}
	return c.rpc.ListenChat(ctx, req, c.callOpts(opts)...)
}

func (c *ListenerClient) ListenCommandBlock(ctx context.Context, req *listenerpb.ListenCommandBlockRequest, opts ...grpc.CallOption) (listenerpb.ListenerService_ListenCommandBlockClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	if req == nil {
		req = &listenerpb.ListenCommandBlockRequest{}
	}
	return c.rpc.ListenCommandBlock(ctx, req, c.callOpts(opts)...)
}
