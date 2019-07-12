package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// Serve starts a given gRPC server.
func Serve(ctx context.Context, server *grpc.Server, address string) (err error) {
	var (
		listen net.Listener
	)

	if listen, err = net.Listen("tcp", address); err != nil {
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

	err = server.Serve(listen)

	return
}
