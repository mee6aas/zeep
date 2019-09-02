package agent

import (
	"context"

	"google.golang.org/grpc"

	server "github.com/mee6aas/zeep/pkg/protocol/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"
	invokerV1Svc "github.com/mee6aas/zeep/pkg/service/invoker/v1"

	invokeeV1Handle "github.com/mee6aas/zeep/internal/pkg/agent/handle/invokee/v1"
	invokerV1Handle "github.com/mee6aas/zeep/internal/pkg/agent/handle/invoker/v1"
)

var (
	gRPCServer *grpc.Server
)

// Serve starts services.
func Serve(ctx context.Context, addr string) (e error) {
	gRPCServer = grpc.NewServer()

	invokeeV1API.RegisterInvokeeServer(gRPCServer, invokeeV1Svc.NewInvokeeAPIServer(
		invokeeV1Handle.Handle{
			WorkerPool: &workerPool,
		},
	))
	invokerV1API.RegisterInvokerServer(gRPCServer, invokerV1Svc.NewInvokerAPIServer(
		invokerV1Handle.Handle{
			WorkerPool: &workerPool,
		},
	))

	if e = server.Serve(ctx, gRPCServer, addr); e != nil {
		return
	}

	return
}
