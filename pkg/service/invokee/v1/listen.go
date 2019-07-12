package v1

import (
	"context"
	"net"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Task is message to be sent to invokee.
type Task = apiV1.Task

// Listen accepts client.
func (s *invokeeAPIServer) Listen(in *apiV1.ListenRequest, stream apiV1.Invokee_ListenServer) (err error) {
	var (
		addr      *net.TCPAddr
		conn      chan Task
		ctxStream context.Context
		ccStream  context.CancelFunc
	)

	if p, ok := peer.FromContext(stream.Context()); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		err = status.Error(codes.Unknown, "Failed to resolve connection information")
		return
	}

	conn = make(chan Task, 1)
	ctxStream, ccStream = context.WithCancel(context.Background())
	defer ccStream()

	if err = s.handle.Connected(ctxStream, addr, conn); err != nil {
		err = status.Error(codes.PermissionDenied, "Operation refused")
		return
	}

	go func() {
		defer ccStream()
		for task := range conn {
			if e := stream.Send(&task); e != nil {
				return
			}
		}
	}()

	<-ctxStream.Done()
	s.handle.Disconnected(addr)

	return
}
