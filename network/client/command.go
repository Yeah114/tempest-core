package client

import (
	"context"
	"strings"

	commandpb "github.com/Yeah114/tempest-core/network_api/command"
	"google.golang.org/grpc"
)

type CommandClient struct {
	rpc         commandpb.CommandServiceClient
	callOptions []grpc.CallOption
}

func newCommandClient(rpc commandpb.CommandServiceClient, callOptions []grpc.CallOption) *CommandClient {
	return &CommandClient{
		rpc:         rpc,
		callOptions: append([]grpc.CallOption(nil), callOptions...),
	}
}

func (c *CommandClient) SendWOCommand(ctx context.Context, cmd string, opts ...grpc.CallOption) error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("command")
	}
	req := &commandpb.SendWOCommandRequest{Cmd: strings.TrimSpace(cmd)}
	resp, err := c.rpc.SendWOCommand(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *CommandClient) SendWSCommand(ctx context.Context, cmd string, opts ...grpc.CallOption) error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("command")
	}
	req := &commandpb.SendWSCommandRequest{Cmd: strings.TrimSpace(cmd)}
	resp, err := c.rpc.SendWSCommand(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *CommandClient) SendPlayerCommand(ctx context.Context, cmd string, opts ...grpc.CallOption) error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("command")
	}
	req := &commandpb.SendPlayerCommandRequest{Cmd: strings.TrimSpace(cmd)}
	resp, err := c.rpc.SendPlayerCommand(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *CommandClient) SendAICommand(ctx context.Context, runtimeID, cmd string, opts ...grpc.CallOption) error {
	if c == nil || c.rpc == nil {
		return clientUnavailable("command")
	}
	req := &commandpb.SendAICommandRequest{
		RuntimeId: strings.TrimSpace(runtimeID),
		Cmd:       strings.TrimSpace(cmd),
	}
	resp, err := c.rpc.SendAICommand(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return err
	}
	_, err = generalPayload(resp)
	return err
}

func (c *CommandClient) SendWSCommandWithResponse(ctx context.Context, cmd string, opts ...grpc.CallOption) (string, error) {
	if c == nil || c.rpc == nil {
		return "", clientUnavailable("command")
	}
	req := &commandpb.SendWSCommandWithResponseRequest{Cmd: strings.TrimSpace(cmd)}
	resp, err := c.rpc.SendWSCommandWithResponse(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *CommandClient) SendPlayerCommandWithResponse(ctx context.Context, cmd string, opts ...grpc.CallOption) (string, error) {
	if c == nil || c.rpc == nil {
		return "", clientUnavailable("command")
	}
	req := &commandpb.SendPlayerCommandWithResponseRequest{Cmd: strings.TrimSpace(cmd)}
	resp, err := c.rpc.SendPlayerCommandWithResponse(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}

func (c *CommandClient) SendAICommandWithResponse(ctx context.Context, runtimeID, cmd string, opts ...grpc.CallOption) (string, error) {
	if c == nil || c.rpc == nil {
		return "", clientUnavailable("command")
	}
	req := &commandpb.SendAICommandWithResponseRequest{
		RuntimeId: strings.TrimSpace(runtimeID),
		Cmd:       strings.TrimSpace(cmd),
	}
	resp, err := c.rpc.SendAICommandWithResponse(ctx, req, mergeCallOptions(c.callOptions, opts)...)
	if err != nil {
		return "", err
	}
	return generalPayload(resp)
}
