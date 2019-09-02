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

func (s *invokerAPIServer) Remove(
	ctx context.Context,
	in *apiV1.RemoveRequest,
) (out *apiV1.RemoveResponse, e error) {
	var (
		addr *net.TCPAddr
	)

	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		e = status.Error(codes.FailedPrecondition, "Failed to resolve request information")
		return
	}

	l := log.WithFields(log.Fields{
		"addr": addr.String(),
		"user": in.GetUsername(),
	})

	l.Info("Activity list requested")

	defer func() {
		if e != nil {
			l.WithError(e).Warn("Activity list refused")
		} else {
			l.Info("Activity listed")
		}
	}()

	if e = s.handle.RemoveRequested(ctx,
		in.GetUsername(),
		in.GetActName(),
	); e != nil {
		// TODO: handle not found

		e = status.Error(codes.Unknown, "Unknown")
		return
	}

	out = &apiV1.RemoveResponse{}

	return
}
