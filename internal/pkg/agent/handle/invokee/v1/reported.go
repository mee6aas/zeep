package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mee6aas/zeep/internal/pkg/agent/assigns"
	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

func (h Handle) Reported(req *v1.ReportRequest) (e error) {
	id := req.GetId()

	if ok := assigns.Report(id, req); !ok {
		e = status.Error(codes.NotFound, "Invocation ID not found")
		return
	}

	return
}
