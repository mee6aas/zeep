package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

func (s *invokerAPIServer) Register(
	ctx context.Context,
	in *apiV1.RegisterRequest,
) (out *apiV1.RegisterResponse, e error) {
	out = &apiV1.RegisterResponse{}

	switch in.GetMethod() {
	case apiV1.RegisterMethod_GLOBAL:
		e = status.Error(codes.Unimplemented, "Unimplemented")
		return
	case apiV1.RegisterMethod_LOCAL:
		e = s.handle.RegisterRequested(ctx, in.GetUsername(), in.GetActName(), in.GetPath())
		return

	case apiV1.RegisterMethod_UNKOWN:
		e = status.Error(codes.InvalidArgument, "RegisterMethod UNKNOWN")
		return
	default:
		e = status.Error(codes.InvalidArgument, "Unrecognized RegisterMethod")
	}

	return
}
