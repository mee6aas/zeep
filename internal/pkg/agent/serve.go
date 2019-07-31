package agent

import (
	"context"

	"google.golang.org/grpc"

	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
	server "github.com/mee6aas/zeep/pkg/protocol/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"
	invokerV1Svc "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

// Serve starts services.
func Serve(ctx context.Context, address string) (err error) {

	// TODO: this config and opts are testing perposes.
	if workerPool, err = pool.NewPool(ctx, pool.Config{
		Images: []string{"runtime-nodejs"},
	},
		pool.WithEachCPU(0),
		pool.WithEachMem(0),
		pool.WithMaxCPU(0),
		pool.WithMaxMem(0),
	); err != nil {
		return
	}

	s := grpc.NewServer()

	invokeeV1API.RegisterInvokeeServer(s, invokeeV1Svc.NewInvokeeAPIServer(
		invokeeV1Handle{},
	))
	invokerV1API.RegisterInvokerServer(s, invokerV1Svc.NewInvokerAPIServer(
		invokerV1Handle{},
	))

	server.Serve(ctx, s, address)

	return
}