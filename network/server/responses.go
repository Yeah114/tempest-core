package server

import (
	"context"
	"errors"

	"github.com/Yeah114/tempest-core/network/app"
	responsepb "github.com/Yeah114/tempest-core/network_api/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func generalSuccess(payload string) *responsepb.GeneralResponse {
	return &responsepb.GeneralResponse{
		Status:  responsepb.GeneralResponse_SUCCESS,
		Payload: payload,
	}
}

func boolSuccess(value bool) *responsepb.GeneralBoolResponse {
	return &responsepb.GeneralBoolResponse{
		Status:  responsepb.GeneralBoolResponse_SUCCESS,
		Payload: value,
	}
}

func int32Success(value int32) *responsepb.GeneralInt32Response {
	return &responsepb.GeneralInt32Response{
		Status:  responsepb.GeneralInt32Response_SUCCESS,
		Payload: value,
	}
}

func int64Success(value int64) *responsepb.GeneralInt64Response {
	return &responsepb.GeneralInt64Response{
		Status:  responsepb.GeneralInt64Response_SUCCESS,
		Payload: value,
	}
}

func uint64Success(value uint64) *responsepb.GeneralUint64Response {
	return &responsepb.GeneralUint64Response{
		Status:  responsepb.GeneralUint64Response_SUCCESS,
		Payload: value,
	}
}

// toStatusError converts domain errors into gRPC status errors.
func toStatusError(err error) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok && st.Code() != codes.OK {
		return err
	}
	return status.Error(statusCodeForError(err), err.Error())
}

func statusCodeForError(err error) codes.Code {
	switch {
	case errors.Is(err, app.ErrNotConnected):
		return codes.FailedPrecondition
	case errors.Is(err, app.ErrAlreadyConnected):
		return codes.AlreadyExists
	case errors.Is(err, app.ErrPlayerUUIDUnknown):
		return codes.FailedPrecondition
	case errors.As(err, new(*notFoundError)):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}

type notFoundError struct {
	msg string
}

func (e *notFoundError) Error() string {
	return e.msg
}

func wrapNotFound(msg string) error {
	return &notFoundError{msg: msg}
}
