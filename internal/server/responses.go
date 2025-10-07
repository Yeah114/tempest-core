package server

import (
	"errors"

	responsepb "github.com/Yeah114/tempest-core/network_api/response"
)

func generalSuccess(payload string) *responsepb.GeneralResponse {
	return &responsepb.GeneralResponse{
		Status:  responsepb.GeneralResponse_SUCCESS,
		Payload: payload,
	}
}

func generalFailure(err error) *responsepb.GeneralResponse {
	return &responsepb.GeneralResponse{
		Status:   responsepb.GeneralResponse_FAILED,
		ErrorMsg: errorMessage(err),
	}
}

func boolResponse(value bool, err error) *responsepb.GeneralBoolResponse {
	if err != nil {
		return &responsepb.GeneralBoolResponse{
			Status:   responsepb.GeneralBoolResponse_FAILED,
			ErrorMsg: errorMessage(err),
		}
	}
	return &responsepb.GeneralBoolResponse{
		Status:  responsepb.GeneralBoolResponse_SUCCESS,
		Payload: value,
	}
}

func int32Response(value int32, err error) *responsepb.GeneralInt32Response {
	if err != nil {
		return &responsepb.GeneralInt32Response{
			Status:   responsepb.GeneralInt32Response_FAILED,
			ErrorMsg: errorMessage(err),
		}
	}
	return &responsepb.GeneralInt32Response{
		Status:  responsepb.GeneralInt32Response_SUCCESS,
		Payload: value,
	}
}

func int64Response(value int64, err error) *responsepb.GeneralInt64Response {
	if err != nil {
		return &responsepb.GeneralInt64Response{
			Status:   responsepb.GeneralInt64Response_FAILED,
			ErrorMsg: errorMessage(err),
		}
	}
	return &responsepb.GeneralInt64Response{
		Status:  responsepb.GeneralInt64Response_SUCCESS,
		Payload: value,
	}
}

func uint64Response(value uint64, err error) *responsepb.GeneralUint64Response {
	if err != nil {
		return &responsepb.GeneralUint64Response{
			Status:   responsepb.GeneralUint64Response_FAILED,
			ErrorMsg: errorMessage(err),
		}
	}
	return &responsepb.GeneralUint64Response{
		Status:  responsepb.GeneralUint64Response_SUCCESS,
		Payload: value,
	}
}

func errorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func wrapNotFound(msg string) error {
	return errors.New(msg)
}
