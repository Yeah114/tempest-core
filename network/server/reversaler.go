package server

import (
	"context"

	"github.com/Yeah114/tempest-core/network/app"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	reversalerpb "github.com/Yeah114/tempest-core/network_api/reversaler"
)

// ReversalerService controls connection lifecycle.
type ReversalerService struct {
	reversalerpb.UnimplementedFateReversalerServiceServer
	state *app.FatalderState
}

// NewReversalerService constructs the lifecycle service.
func NewReversalerService(state *app.FatalderState) *ReversalerService {
	return &ReversalerService{state: state}
}

func (s *ReversalerService) NewFateReversaler(ctx context.Context, req *reversalerpb.NewFateReversalerRequest) (*responsepb.GeneralResponse, error) {
	opts := app.ConnectOptions{
		AuthServerAddress: req.GetAuthServer(),
		AuthUsername:      req.GetUserName(),
		AuthPassword:      req.GetUserPassword(),
		AuthToken:         req.GetUserToken(),
		ServerCode:        req.GetServerCode(),
		ServerPassword:    req.GetServerPassword(),
	}
	if err := s.state.Connect(ctx, opts); err != nil {
		return generalFailure(err), nil
	}
	// Warm player registry.
	_, _ = s.state.SnapshotPlayers()
	return generalSuccess(""), nil
}

func (s *ReversalerService) WaitDead(req *reversalerpb.WaitDeadRequest, stream reversalerpb.FateReversalerService_WaitDeadServer) error {
	for {
		ch, cancel := s.state.DisconnectEvents(1)
		select {
		case <-stream.Context().Done():
			cancel()
			return nil
		case err, ok := <-ch:
			cancel()
			if !ok {
				continue
			}
			reason := ""
			if err != nil {
				reason = err.Error()
			}
			if sendErr := stream.Send(&responsepb.DeadReason{Reason: reason}); sendErr != nil {
				return sendErr
			}
			return nil
		}
	}
}

func (s *ReversalerService) Ping(ctx context.Context, req *reversalerpb.PingRequest) (*responsepb.PingResponse, error) {
	return &responsepb.PingResponse{Success: true}, nil
}
