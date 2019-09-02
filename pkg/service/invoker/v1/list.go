package v1

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/mee6aas/zeep/pkg/activity"
	apiV1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

func (s *invokerAPIServer) List(
	ctx context.Context,
	in *apiV1.ListRequest,
) (out *apiV1.ListResponse, e error) {
	var (
		addr *net.TCPAddr
		as   []activity.Activity
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

	if as, e = s.handle.ListRequested(ctx, in.GetUsername()); e != nil {
		// TODO: handle not found

		e = status.Error(codes.Unknown, "Unknown")
		return
	}

	out = &apiV1.ListResponse{}
	out.Activities = make([]*apiV1.ManagedActivity, len(as))

	for i, a := range as {
		out.Activities[i] = &apiV1.ManagedActivity{
			Name:    a.Name,
			Runtime: a.Runtime,
			Added:   a.AddedDate,
		}
	}

	return
}
