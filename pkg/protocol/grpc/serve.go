package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// Serve starts a given gRPC server.
func Serve(ctx context.Context, server *grpc.Server, address string) (e error) {
	var (
		l net.Listener
	)

	if l, e = net.Listen("tcp", address); e != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
		case <-ctx.Done():
		}

		server.GracefulStop()

		return
	}()

	e = server.Serve(l)

	return
}
