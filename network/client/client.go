package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	commandpb "github.com/Yeah114/tempest-core/network_api/command"
	listenerpb "github.com/Yeah114/tempest-core/network_api/listener"
	playerkitpb "github.com/Yeah114/tempest-core/network_api/playerkit"
	reversalerpb "github.com/Yeah114/tempest-core/network_api/reversaler"
	utilspb "github.com/Yeah114/tempest-core/network_api/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ErrTargetRequired indicates that a dial target is mandatory.
var ErrTargetRequired = errors.New("dial target required")

// Option applies custom behaviour when establishing the client connection.
type Option func(*dialConfig)

type dialConfig struct {
	transportCreds credentials.TransportCredentials
	dialOptions    []grpc.DialOption
	callOptions    []grpc.CallOption
}

// WithTransportCredentials overrides the default insecure transport credentials.
func WithTransportCredentials(creds credentials.TransportCredentials) Option {
	return func(cfg *dialConfig) {
		cfg.transportCreds = creds
	}
}

// WithDialOptions appends additional grpc.DialOption values.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(cfg *dialConfig) {
		cfg.dialOptions = append(cfg.dialOptions, opts...)
	}
}

// WithCallOptions appends default grpc.CallOption values used for every RPC.
func WithCallOptions(opts ...grpc.CallOption) Option {
	return func(cfg *dialConfig) {
		cfg.callOptions = append(cfg.callOptions, opts...)
	}
}

// Client aggregates typed tempest-core gRPC clients.
type Client struct {
	conn        *grpc.ClientConn
	callOptions []grpc.CallOption

	Command    *CommandClient
	Listener   *ListenerClient
	PlayerKit  *PlayerKitClient
	Reversaler *ReversalerClient
	Utils      *UtilsClient
}

// Dial connects to a tempest-core gRPC endpoint and initialises service clients.
func Dial(ctx context.Context, target string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(target) == "" {
		return nil, ErrTargetRequired
	}

	cfg := dialConfig{
		transportCreds: insecure.NewCredentials(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	dialOptions := make([]grpc.DialOption, 0, len(cfg.dialOptions)+1)
	if cfg.transportCreds != nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(cfg.transportCreds))
	}
	dialOptions = append(dialOptions, cfg.dialOptions...)

	conn, err := grpc.DialContext(ctx, target, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("tempest client dial %q: %w", target, err)
	}

	callOpts := append([]grpc.CallOption(nil), cfg.callOptions...)
	c := &Client{
		conn:        conn,
		callOptions: callOpts,
	}
	c.Command = newCommandClient(commandpb.NewCommandServiceClient(conn), callOpts)
	c.Listener = newListenerClient(listenerpb.NewListenerServiceClient(conn), callOpts)
	c.PlayerKit = newPlayerKitClient(playerkitpb.NewPlayerKitServiceClient(conn), callOpts)
	c.Reversaler = newReversalerClient(reversalerpb.NewFateReversalerServiceClient(conn), callOpts)
	c.Utils = newUtilsClient(utilspb.NewUtilsServiceClient(conn), callOpts)

	return c, nil
}

// Conn exposes the underlying grpc.ClientConn.
func (c *Client) Conn() *grpc.ClientConn {
	if c == nil {
		return nil
	}
	return c.conn
}

// Close terminates the underlying gRPC connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
