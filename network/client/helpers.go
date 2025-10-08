package client

import (
	"fmt"

	"google.golang.org/grpc"
)

func mergeCallOptions(base []grpc.CallOption, opts []grpc.CallOption) []grpc.CallOption {
	if len(base) == 0 {
		return append([]grpc.CallOption(nil), opts...)
	}
	out := make([]grpc.CallOption, 0, len(base)+len(opts))
	out = append(out, base...)
	out = append(out, opts...)
	return out
}

func clientUnavailable(name string) error {
	return fmt.Errorf("%s client not initialised", name)
}
