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

// Handover handles listener handover.
func (s *invokeeAPIServer) Handover(
	ctx context.Context,
	in *apiV1.HandoverRequest,
) (out *apiV1.HandoverResponse, e error) {
	var (
		addr *net.TCPAddr
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

	l.Info("Worker handover requested")

	if e = s.handle.HandoverRequested(addr); e != nil {
		e = status.Error(codes.PermissionDenied, "Non-listener worker has no control to handover")
		return
	}

	out = &apiV1.HandoverResponse{}

	return
}
