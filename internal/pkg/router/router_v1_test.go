package router_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/router"
	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockHandle struct {
	refucseConnection bool
	isConnected       bool
}

func (h *mockHandle) Connected(*net.TCPAddr) (err error) {
	h.isConnected = true

	if h.refucseConnection {
		err = errors.New("refuse")
	}
	return
}
func (h *mockHandle) Disconnected(*net.TCPAddr) { h.isConnected = false }
func (h *mockHandle) SendRequested(_ *v1.SendRequest, addr *net.TCPAddr) (*v1.Message, *net.TCPAddr, error) {
	return &v1.Message{}, addr, nil
}

var (
	testServeFailed  = false
	testListenFailed = false

	testWaitTime = time.Millisecond * 100
	testHandle   = &mockHandle{
		refucseConnection: false,
		isConnected:       false,
	}
	testRouter    = router.NewRouterV1(testHandle)
	testServerFin = make(chan struct{})
	testStreamFin = make(chan struct{})
)

func TestServe(t *testing.T) {
	var (
		err error
	)

	go func() {
		err = testRouter.Serve()
		testServerFin <- struct{}{}
	}()

	time.Sleep(testWaitTime)

	if err != nil {
		testServeFailed = true
		t.Fatalf("failed to serve router %v", err)
	}

	testRouter.Stop()

	<-testServerFin
}

func TestListen(t *testing.T) {
	var (
		err    error
		conn   *grpc.ClientConn
		stream v1.RouterService_ListenClient
	)

	if testServeFailed {
		t.Skipf("TestServe failed")
	}

	// keep server running.
	// this server used at later test.
	go func() {
		err = testRouter.Serve()
		testServerFin <- struct{}{}
	}()

	time.Sleep(testWaitTime)

	conn, err = grpc.Dial("localhost:5122", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial %v", err)
	}

	defer conn.Close()

	client := v1.NewRouterServiceClient(conn)

	// test handle.Connected and handle.Disconnect
	{
		ctx, cancel := context.WithTimeout(context.Background(), testWaitTime*5)

		_, err = client.Listen(ctx, &v1.ListenRequest{})
		if err != nil {
			testListenFailed = true
			t.Fatalf("failed to listen request %v", err)
		}

		time.Sleep(testWaitTime)

		if !testHandle.isConnected {
			testListenFailed = true
			t.Fatal("handle.Connected does not seem to work")
		}

		cancel()
		time.Sleep(testWaitTime)

		if testHandle.isConnected {
			testListenFailed = true
			t.Fatal("handle.Disconnected does not seem to work")
		}
	}

	// test connection refuse
	{
		testHandle.refucseConnection = true
		ctx, cancel := context.WithTimeout(context.Background(), testWaitTime*5)

		stream, err = client.Listen(ctx, &v1.ListenRequest{})
		if err != nil {
			testListenFailed = true
			t.Fatalf("failed to listen request %v", err)
		}

		time.Sleep(testWaitTime)

		if !testHandle.isConnected {
			testListenFailed = true
			t.Fatal("handle does not seem to work")
		}

		_, err = stream.Recv()

		s, ok := status.FromError(err)

		if !ok {
			testListenFailed = true
			t.Fatalf("expected connection be refused %v", err)
		}

		switch s.Code() {
		case codes.PermissionDenied:
		default:
			testListenFailed = true
			t.Fatalf("expected error be PermissionDenied %v", err)
		}

		cancel()
		testHandle.refucseConnection = false
	}
}

func TestSend(t *testing.T) {
	var (
		err    error
		conn   *grpc.ClientConn
		stream v1.RouterService_ListenClient
	)

	if testListenFailed {
		t.Skipf("TestServe failed")
	}

	conn, err = grpc.Dial("localhost:5122", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial %v", err)
	}

	defer conn.Close()

	client := v1.NewRouterServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), testWaitTime*5)
	defer cancel()

	stream, err = client.Listen(ctx, &v1.ListenRequest{})
	if err != nil {
		testListenFailed = true
		t.Fatalf("failed to listen request %v", err)
	}

	_, err = client.Send(ctx, &v1.SendRequest{})
	if err != nil {
		t.Fatalf("failed to send request %v", err)
	}

	_, err = stream.Recv()
	if err != nil {
		t.Fatalf("failed to receive %v", err)
	}
}
