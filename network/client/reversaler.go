package client

import (
	"context"
	"strings"

	reversalerpb "github.com/Yeah114/tempest-core/network_api/reversaler"
	"google.golang.org/grpc"
)

type ConnectOptions struct {
	AuthServer     string
	Username       string
	Password       string
	Token          string
	ServerCode     string
	ServerPassword string
}

type ReversalerClient struct {
	rpc         reversalerpb.FateReversalerServiceClient
	callOptions []grpc.CallOption
}

func newReversalerClient(rpc reversalerpb.FateReversalerServiceClient, callOptions []grpc.CallOption) *ReversalerClient {
	return &ReversalerClient{
		rpc:         rpc,
		callOptions: append([]grpc.CallOption(nil), callOptions...),
	}
}

func (c *ReversalerClient) ready() error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("reversaler")
	}
	return nil
}

func (c *ReversalerClient) callOpts(opts []grpc.CallOption) []grpc.CallOption {
	return mergeCallOptions(c.callOptions, opts)
}

func (c *ReversalerClient) Connect(ctx context.Context, options ConnectOptions, opts ...grpc.CallOption) error {
	if err := c.ready(); err != nil {
		return err
	}
	req := &reversalerpb.NewFateReversalerRequest{
		AuthServer:     strings.TrimSpace(options.AuthServer),
		UserName:       strings.TrimSpace(options.Username),
		UserPassword:   options.Password,
		UserToken:      strings.TrimSpace(options.Token),
		ServerCode:     strings.TrimSpace(options.ServerCode),
		ServerPassword: strings.TrimSpace(options.ServerPassword),
	}
	resp, err := c.rpc.NewFateReversaler(ctx, req, c.callOpts(opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *ReversalerClient) Ping(ctx context.Context, opts ...grpc.CallOption) (bool, error) {
	if err := c.ready(); err != nil {
		return false, err
	}
	resp, err := c.rpc.Ping(ctx, &reversalerpb.PingRequest{}, c.callOpts(opts)...)
	if err != nil {
		return false, err
	}
	return resp.GetSuccess(), nil
}

func (c *ReversalerClient) WaitDead(ctx context.Context, opts ...grpc.CallOption) (reversalerpb.FateReversalerService_WaitDeadClient, error) {
	if err := c.ready(); err != nil {
		return nil, err
	}
	return c.rpc.WaitDead(ctx, &reversalerpb.WaitDeadRequest{}, c.callOpts(opts)...)
}
