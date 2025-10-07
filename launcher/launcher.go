package launcher

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Yeah114/tempest-core/internal/app"
	core "github.com/Yeah114/tempest-core/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Options defines the listening configuration for the embedded gRPC server.
type Options struct {
	Address  string
	Port     int
	Callback func()
}

// Server manages the lifecycle of a tempest-core gRPC server.
type Server struct {
	opts  Options
	srv   *grpc.Server
	lis   net.Listener
	state *app.FatalderState

	once sync.Once
}

// Start launches a tempest-core server bound to the provided address.
func Start(ctx context.Context, opts Options) (*Server, error) {
	if opts.Address == "" {
		opts.Address = "0.0.0.0"
	}
	if opts.Port == 0 {
		opts.Port = 20919
	}

	addr := fmt.Sprintf("%s:%d", opts.Address, opts.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("launcher: listen %s: %w", addr, err)
	}

	state := app.NewFatalderState()
	services := core.NewServices(state)

	srv := grpc.NewServer()
	services.Register(srv)
	reflection.Register(srv)

	l := &Server{
		opts:  opts,
		srv:   srv,
		lis:   lis,
		state: state,
	}

	go func() {
		<-ctx.Done()
		l.Stop()
	}()

	go func() {
		if err := srv.Serve(lis); err != nil {
			l.Stop()
		}
		if opts.Callback != nil {
			opts.Callback()
		}
	}()

	return l, nil
}

// Address returns the listener address.
func (s *Server) Address() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.opts.Address, s.opts.Port)
}

// Stop gracefully shuts down the gRPC server and active connection.
func (s *Server) Stop() {
	if s == nil {
		return
	}
	s.once.Do(func() {
		if s.srv != nil {
			s.srv.GracefulStop()
		}
		if s.lis != nil {
			_ = s.lis.Close()
		}
		if s.state != nil {
			_ = s.state.Disconnect()
		}
	})
}
