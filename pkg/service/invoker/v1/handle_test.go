package v1_test

import (
	v1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type mockHandle struct {
}

func (h *mockHandle) Requested(
	in *v1.InvokeRequest,
) (out *v1.InvokeResponse, err error) {
	act := in.GetTarget()
	out = &v1.InvokeResponse{Result: act.GetName() + " in " + act.GetId()}
	return
}
