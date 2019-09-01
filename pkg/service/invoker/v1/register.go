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

func (s *invokerAPIServer) Register(
	ctx context.Context,
	in *apiV1.RegisterRequest,
) (out *apiV1.RegisterResponse, e error) {
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
		"user": in.GetUsername(),
		"name": in.GetActName(),
		"path": in.GetPath(),
	})

	l.Info("Activity register requested")

	out = &apiV1.RegisterResponse{}

	switch in.GetMethod() {
	case apiV1.RegisterMethod_GLOBAL:
		e = status.Error(codes.Unimplemented, "Unimplemented")
	case apiV1.RegisterMethod_LOCAL:
		e = s.handle.RegisterRequested(ctx, in.GetUsername(), in.GetActName(), in.GetPath())
	case apiV1.RegisterMethod_UNKOWN:
		e = status.Error(codes.InvalidArgument, "RegisterMethod UNKNOWN")
	default:
		e = status.Error(codes.InvalidArgument, "Unrecognized RegisterMethod")
	}

	if e != nil {
		l.WithError(e).Warn("Activity register refused")
	} else {
		l.Info("Activity registered")
	}

	return
}
