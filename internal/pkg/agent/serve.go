package agent

import (
	"context"

	"google.golang.org/grpc"

	server "github.com/mee6aas/zeep/pkg/protocol/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"
	invokerV1Svc "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

// Serve starts services.
func Serve(ctx context.Context, address string) (e error) {
	s := grpc.NewServer()

	invokeeV1API.RegisterInvokeeServer(s, invokeeV1Svc.NewInvokeeAPIServer(
		invokeeV1Handle{},
	))
	invokerV1API.RegisterInvokerServer(s, invokerV1Svc.NewInvokerAPIServer(
		invokerV1Handle{},
	))

	if e = server.Serve(ctx, s, address); e != nil {
		return
	}

	return
}
