package v1

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Task is message to be sent to invokee.
type Task = apiV1.Task

// Listen accepts client.
func (s *invokeeAPIServer) Listen(
	in *apiV1.ListenRequest,
	stream apiV1.Invokee_ListenServer,
) (e error) {
	var (
		addr      *net.TCPAddr
		conn      chan Task
		ctxStream context.Context
		ccStream  context.CancelFunc
	)

	if p, ok := peer.FromContext(stream.Context()); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		e = status.Error(codes.Unknown, "Failed to resolve request information")
		return
	}

	l := log.WithFields(log.Fields{
		"addr": addr.String(),
	})

	l.Info("Worker listen requested")

	defer func() {
		if e != nil {
			l.WithError(e).Warn("Worker listen refused")
		}
	}()

	conn = make(chan Task, 1)
	ctxStream, ccStream = context.WithCancel(context.Background())
	defer ccStream()

	if e = s.handle.Connected(ctxStream, addr, conn); e != nil {
		e = status.Errorf(codes.PermissionDenied, "Operation refused: %s", e.Error())
		return
	}

	l.Info("Connected")

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

	l.Info("Disconnected")

	return
}
