package v1

import (
	"context"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// InvocationResult is result of invocation from the invokee.
type InvocationResult = apiV1.ReportRequest

// Report handles result of invocation.
func (s *invokeeAPIServer) Report(ctx context.Context, in *apiV1.ReportRequest) (rst *apiV1.ReportResponse, err error) {
	rst = &apiV1.ReportResponse{}

	s.handle.Reported(in)

	return
}
