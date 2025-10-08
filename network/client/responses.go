package client

import (
	"errors"
	"strings"

	responsepb "github.com/Yeah114/tempest-core/network_api/response"
)

var (
	errNilResponse   = errors.New("nil response")
	errUnknownStatus = errors.New("remote operation failed")
)

func generalPayload(resp *responsepb.GeneralResponse) (string, error) {
	if resp == nil {
		return "", errNilResponse
	}
	if resp.GetStatus() != responsepb.GeneralResponse_SUCCESS {
		return "", responseError(resp.GetErrorMsg())
	}
	return resp.GetPayload(), nil
}

func boolPayload(resp *responsepb.GeneralBoolResponse) (bool, error) {
	if resp == nil {
		return false, errNilResponse
	}
	if resp.GetStatus() != responsepb.GeneralBoolResponse_SUCCESS {
		return false, responseError(resp.GetErrorMsg())
	}
	return resp.GetPayload(), nil
}

func int32Payload(resp *responsepb.GeneralInt32Response) (int32, error) {
	if resp == nil {
		return 0, errNilResponse
	}
	if resp.GetStatus() != responsepb.GeneralInt32Response_SUCCESS {
		return 0, responseError(resp.GetErrorMsg())
	}
	return resp.GetPayload(), nil
}

func int64Payload(resp *responsepb.GeneralInt64Response) (int64, error) {
	if resp == nil {
		return 0, errNilResponse
	}
	if resp.GetStatus() != responsepb.GeneralInt64Response_SUCCESS {
		return 0, responseError(resp.GetErrorMsg())
	}
	return resp.GetPayload(), nil
}

func uint64Payload(resp *responsepb.GeneralUint64Response) (uint64, error) {
	if resp == nil {
		return 0, errNilResponse
	}
	if resp.GetStatus() != responsepb.GeneralUint64Response_SUCCESS {
		return 0, responseError(resp.GetErrorMsg())
	}
	return resp.GetPayload(), nil
}

func responseError(msg string) error {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return errUnknownStatus
	}
	return errors.New(msg)
}
