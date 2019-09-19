package v1

import (
	"net"

	"github.com/pkg/errors"
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
		ctx  = stream.Context()
		err  error
		addr *net.TCPAddr
		conn chan Task
	)

	if p, ok := peer.FromContext(ctx); ok {
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
		if err != nil {
			l.WithError(err).Warn("Worker listen failed with error")
		}
	}()

	conn = make(chan Task, 1)

	if err = s.handle.Connected(ctx, addr, conn); err != nil {
		err = errors.Wrap(err, "Handle:Connected returns error")
		e = status.Errorf(codes.PermissionDenied, "Operation refused: %s", err.Error())

		return
	}

	l.Info("Connected")

FORWARD:
	for {
		select {
		case task, ok := <-conn:
			if !ok {
				// Task assigner closes the channel
				break FORWARD
			}

			if err = stream.Send(&task); err != nil {
				err = errors.Wrap(err, "Failed to send task")

				break
			}

		case <-ctx.Done():
			break FORWARD
		}
	}

	s.handle.Disconnected(addr)

	l.Info("Disconnected")

	return
}
