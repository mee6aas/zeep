package router

import (
	"log"
	"net"

	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Listen accepts client.
func (r *V1) Listen(in *v1.ListenRequest, stream v1.RouterService_ListenServer) (err error) {
	var (
		addr *net.TCPAddr
		conn chan v1.Message
	)

	if p, ok := peer.FromContext(stream.Context()); ok {
		addr = p.Addr.(*net.TCPAddr)
	} else {
		err = status.Error(codes.Unknown, "failed to resolve connection information")
		return
	}

	conn = make(chan v1.Message, 5)
	connections[addr.IP.String()] = conn

	err = r.handle.Connected(addr)
	if err != nil {
		delete(connections, addr.IP.String())
		close(conn)

		err = status.Error(codes.PermissionDenied, "operation refused")
		return
	}

	go func() {
		for {
			select {
			case <-stream.Context().Done():
				return
			case msg := <-conn:
				err = stream.Send(&msg)
				if err != nil {
					log.Printf("is stream closed? %v", err)
					return
				}
			}
		}
	}()

	<-stream.Context().Done()
	r.handle.Disconnected(addr)
	delete(connections, addr.IP.String())
	close(conn)

	return
}
