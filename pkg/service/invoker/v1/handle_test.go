package v1_test

import (
	"context"
	"fmt"

	"github.com/mee6aas/zeep/pkg/activity"
	v1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type mockHandle struct {
}

func (h *mockHandle) InvokeRequested(
	_ context.Context,
	username string,
	actName string,
	arg string,
) (out *v1.InvokeResponse, e error) {
	out = &v1.InvokeResponse{Result: fmt.Sprintf("%s by %s with %s", actName, username, arg)}
	return
}

func (h *mockHandle) AddRequested(
	_ context.Context,
	username string,
	actName string,
	actDir string,
) (e error) {
	return
}

func (h *mockHandle) ListRequested(
	_ context.Context,
	username string,
) (out []activity.Activity, e error) {
	return
}

func (h *mockHandle) RemoveRequested(
	_ context.Context,
	username string,
	actName string,
) (e error) {
	return
}
