package server

import (
	"github.com/Yeah114/tempest-core/network/app"
	commandpb "github.com/Yeah114/tempest-core/network_api/command"
	listenerpb "github.com/Yeah114/tempest-core/network_api/listener"
	playerkitpb "github.com/Yeah114/tempest-core/network_api/playerkit"
	reversalerpb "github.com/Yeah114/tempest-core/network_api/reversaler"
	utilspb "github.com/Yeah114/tempest-core/network_api/utils"
	"google.golang.org/grpc"
)

// Services bundles all gRPC handlers.
type Services struct {
	Command    *CommandService
	Listener   *ListenerService
	PlayerKit  *PlayerKitService
	Reversaler *ReversalerService
	Utils      *UtilsService
}

// NewServices wires up every service against shared state.
func NewServices(state *app.FatalderState) *Services {
	return &Services{
		Command:    NewCommandService(state),
		Listener:   NewListenerService(state),
		PlayerKit:  NewPlayerKitService(state),
		Reversaler: NewReversalerService(state),
		Utils:      NewUtilsService(state),
	}
}

// Register attaches all services to the provided gRPC server.
func (s *Services) Register(server *grpc.Server) {
	if s == nil {
		return
	}
	commandpb.RegisterCommandServiceServer(server, s.Command)
	listenerpb.RegisterListenerServiceServer(server, s.Listener)
	playerkitpb.RegisterPlayerKitServiceServer(server, s.PlayerKit)
	reversalerpb.RegisterFateReversalerServiceServer(server, s.Reversaler)
	utilspb.RegisterUtilsServiceServer(server, s.Utils)
}
