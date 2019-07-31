package v1

import (
	"context"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// InvokeRequest is a request to invoke an activity.
type InvokeRequest = apiV1.InvokeRequest

// InvokeResponse is a response of invocation.
type InvokeResponse = apiV1.InvokeResponse

func (s *invokerAPIServer) Invoke(
	ctx context.Context,
	in *apiV1.InvokeRequest,
) (out *apiV1.InvokeResponse, err error) {
	out, err = s.handle.Requested(ctx, in)
	return
}