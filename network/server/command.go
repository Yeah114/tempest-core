package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	fpacket "github.com/Yeah114/FunShuttler/core/minecraft/protocol/packet"
	"github.com/Yeah114/FunShuttler/game_control/game_interface"
	"github.com/Yeah114/tempest-core/network/app"
	commandpb "github.com/Yeah114/tempest-core/network_api/command"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CommandService implements gRPC command endpoints.
type CommandService struct {
	commandpb.UnimplementedCommandServiceServer
	state *app.FatalderState
}

// NewCommandService constructs a command service bound to shared state.
func NewCommandService(state *app.FatalderState) *CommandService {
	return &CommandService{state: state}
}

func (s *CommandService) SendWOCommand(ctx context.Context, req *commandpb.SendWOCommandRequest) (*responsepb.GeneralResponse, error) {
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		return iface.Commands().SendSettingsCommand(strings.TrimSpace(req.GetCmd()), false)
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *CommandService) SendWSCommand(ctx context.Context, req *commandpb.SendWSCommandRequest) (*responsepb.GeneralResponse, error) {
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		return iface.Commands().SendWSCommand(strings.TrimSpace(req.GetCmd()))
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *CommandService) SendPlayerCommand(ctx context.Context, req *commandpb.SendPlayerCommandRequest) (*responsepb.GeneralResponse, error) {
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		return iface.Commands().SendPlayerCommand(strings.TrimSpace(req.GetCmd()))
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *CommandService) SendAICommand(ctx context.Context, req *commandpb.SendAICommandRequest) (*responsepb.GeneralResponse, error) {
	cmd := buildAIExecute(strings.TrimSpace(req.GetRuntimeId()), strings.TrimSpace(req.GetCmd()))
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		return iface.Commands().SendWSCommand(cmd)
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(""), nil
}

func (s *CommandService) SendWSCommandWithResponse(ctx context.Context, req *commandpb.SendWSCommandWithResponseRequest) (*responsepb.GeneralResponse, error) {
	var output *fpacket.CommandOutput
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		var err error
		output, err = iface.Commands().SendWSCommandWithResp(strings.TrimSpace(req.GetCmd()))
		return err
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	payload, err := marshalCommandOutput(output)
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(payload), nil
}

func (s *CommandService) SendPlayerCommandWithResponse(ctx context.Context, req *commandpb.SendPlayerCommandWithResponseRequest) (*responsepb.GeneralResponse, error) {
	var output *fpacket.CommandOutput
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		var err error
		output, err = iface.Commands().SendPlayerCommandWithResp(strings.TrimSpace(req.GetCmd()))
		return err
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	payload, err := marshalCommandOutput(output)
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(payload), nil
}

func (s *CommandService) SendAICommandWithResponse(ctx context.Context, req *commandpb.SendAICommandWithResponseRequest) (*responsepb.GeneralResponse, error) {
	runtimeID := strings.TrimSpace(req.GetRuntimeId())
	if runtimeID == "" {
		return nil, status.Error(codes.InvalidArgument, "runtime_id required")
	}
	var output *fpacket.CommandOutput
	cmd := buildAIExecute(runtimeID, strings.TrimSpace(req.GetCmd()))
	err := s.state.WithGameInterface(func(iface *game_interface.GameInterface) error {
		var err error
		output, err = iface.Commands().SendWSCommandWithResp(cmd)
		return err
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	payload, err := marshalCommandOutput(output)
	if err != nil {
		return nil, toStatusError(err)
	}
	return generalSuccess(payload), nil
}

func marshalCommandOutput(output *fpacket.CommandOutput) (string, error) {
	if output == nil {
		return "{}", nil
	}
	bs, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func buildAIExecute(runtimeID, cmd string) string {
	if runtimeID == "" {
		return cmd
	}
	if cmd == "" {
		return fmt.Sprintf("execute as @e[runtime_id=%s] run ", runtimeID)
	}
	return fmt.Sprintf("execute as @e[runtime_id=%s] at @s run %s", runtimeID, cmd)
}
