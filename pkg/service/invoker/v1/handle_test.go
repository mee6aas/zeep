package v1_test

import (
	"context"
	"fmt"

	v1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type mockHandle struct {
}

func (h *mockHandle) InvokeRequested(
	_ context.Context,
	username string,
	actName string,
	actLabel string,
	arg string,
) (out *v1.InvokeResponse, err error) {
	out = &v1.InvokeResponse{Result: fmt.Sprintf("%s@%s by %s with %s", actName, actLabel, username, arg)}
	return
}

func (h *mockHandle) RegisterRequested(
	_ context.Context,
	username string,
	actDir string,
) (err error) {
	return
}
