package agent

import (
	"context"

	"google.golang.org/grpc"

	invokeeV1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	server "github.com/mee6aas/zeep/pkg/protocol/grpc"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"
)

// Serve starts services.
func Serve(ctx context.Context, address string) (err error) {

	s := grpc.NewServer()

	invokeeV1.RegisterInvokeeServer(s, invokeeV1Svc.NewInvokeeAPIServer(
		invokeeV1Handle{},
	))
	// invokerV1.RegisterInvokerServer(s, a)

	server.Serve(ctx, s, address)

	return
}
