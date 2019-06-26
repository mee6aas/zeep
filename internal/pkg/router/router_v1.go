package router

import (
	"net"

	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"
	"google.golang.org/grpc"
)

// V1 _
type V1 struct {
	handle HandleV1
	server *grpc.Server
}

// NewRouterV1 creates RouterV1.
func NewRouterV1(handle HandleV1) (r V1) {
	r.handle = handle

	return
}

// Serve strats grpc server.
func (r *V1) Serve() (err error) {
	var (
		l net.Listener
		s *grpc.Server
	)

	l, err = net.Listen("tcp", "localhost:5122")
	if err != nil {
		return
	}

	s = grpc.NewServer()
	v1.RegisterRouterServiceServer(s, r)
	r.server = s

	err = s.Serve(l)
	if err != nil {
		return
	}

	return
}

// Stop stops grpc server.
func (r *V1) Stop() {
	r.server.Stop()
}
