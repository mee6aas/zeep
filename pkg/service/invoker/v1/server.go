package v1

import (
	apiV1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

type invokerAPIServer struct {
	handle InvokerAPIServerHandle
}

// NewInvokerAPIServer creates API server for Invoker service.
func NewInvokerAPIServer(h InvokerAPIServerHandle) apiV1.InvokerServer {
	return &invokerAPIServer{handle: h}
}
