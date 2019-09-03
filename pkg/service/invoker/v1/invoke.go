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
		if e != nil {
			l.WithError(e).Warn("Activity invoke refused")
		} else {
			l.Info("Activity invoke added")
		}
	}()

	if username == "" {
		if username, e = s.handle.ResolveNameFromIP(ctx, addr.IP.String()); e != nil {
			// username not found
			e = status.Error(codes.PermissionDenied, "Username not found from IP")
			return
		}
	}

	if out, e = s.handle.InvokeRequested(ctx,
		username,
		in.GetActName(),
		in.GetArg(),
	); e != nil {
		e = status.Error(codes.Unknown, "Unknown")
		return
	}

	return
}
