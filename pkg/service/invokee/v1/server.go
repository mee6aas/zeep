package v1

import (
	apiV1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

type invokeeAPIServer struct {
	handle InvokeeAPIServerHandle
}

// NewInovkeeAPIServer creates API server for Invokee service.
func NewInovkeeAPIServer(h InvokeeAPIServerHandle) apiV1.InvokeeServer {
	return &invokeeAPIServer{handle: h}
}
