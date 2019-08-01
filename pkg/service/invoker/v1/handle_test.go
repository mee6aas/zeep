package v1_test

import (
	"context"

	v1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type mockHandle struct {
}

func (h *mockHandle) InvokeRequested(
	_ context.Context,
	in *v1.InvokeRequest,
) (out *v1.InvokeResponse, err error) {
	act := in.GetTarget()
	out = &v1.InvokeResponse{Result: act.GetName() + "@" + act.GetLabel()}
	return
}

func (h *mockHandle) RegisterRequested(
	_ context.Context,
	username string,
	actDir string,
) (err error) {
	return
}
