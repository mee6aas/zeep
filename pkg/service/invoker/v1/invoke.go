package v1

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	apiV1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// InvokeRequest is a request to invoke an activity.
type InvokeRequest = apiV1.InvokeRequest

// InvokeResponse is a response of invocation.
type InvokeResponse = apiV1.InvokeResponse

func (s *invokerAPIServer) Invoke(
	ctx context.Context,
	in *apiV1.InvokeRequest,
) (out *apiV1.InvokeResponse, e error) {
	var (
		err      error
		addr     *net.TCPAddr
		username string
	)

	username = in.GetUsername()

	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		e = status.Error(codes.Unknown, "Failed to resolve request information")
		return
	}

	l := log.WithFields(log.Fields{
		"addr": addr.String(),
		"user": username,
		"name": in.GetActName(),
	})

	l.Info("Activity invoke requested")

	defer func() {
		if err != nil {
			l.WithError(err).Warn("Activity invoke refused")
		} else {
			l.Info("Activity invoked")
		}
	}()

	if out, err = s.handle.InvokeRequested(ctx,
		addr,
		username,
		in.GetActName(),
		in.GetArg(),
	); err != nil {
		e = status.Error(codes.Unknown, "Unknown")
		return
	}

	return
}
