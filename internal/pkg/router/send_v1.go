package router

import (
	"context"
	"net"

	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Send sends message to given destination.
func (r *V1) Send(ctx context.Context, in *v1.SendRequest) (out *v1.SendResponse, err error) {
	var (
		addr *net.TCPAddr
		msg  *v1.Message
		conn chan v1.Message
	)

	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		err = status.Error(codes.Unknown, "failed to resolve connection information")
		return
	}

	msg, addr, err = r.handle.SendRequested(in, addr)
	if err != nil {
		err = status.Errorf(codes.PermissionDenied, "operation refused")
		return
	}

	if c, ok := connections[addr.IP.String()]; ok {
		conn = c
	} else {
		err = status.Error(codes.Internal, "invalid forwarding by handle")
		return
	}

	conn <- *msg

	out = &v1.SendResponse{}

	return
}
