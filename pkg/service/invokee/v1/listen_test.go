package v1_test

import (
	"context"
	"log"
	"testing"
	"time"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

func TestListen(t *testing.T) {
	var (
		err    error
		ctx    context.Context
		cc     context.CancelFunc
		stream v1.Invokee_ListenClient
	)

	ctx, cc = context.WithTimeout(context.Background(), time.Second)
	defer cc()

	if stream, err = testClient.Listen(ctx, &v1.ListenRequest{}); err != nil {
		t.Fatalf("failed to request to listen: %v", err)
	}

	for {
		if _, err = stream.Recv(); err != nil {
			break
		}

		log.Println("data received")
	}

	select {
	case <-stream.Context().Done():
	case <-time.After(time.Millisecond * 500):
		t.Fatal("failed to close")
	}

	time.Sleep(time.Millisecond * 100)
}

func TestDisconnectWhileListening(t *testing.T) {
	var (
		err    error
		ctx    context.Context
		cc     context.CancelFunc
		stream v1.Invokee_ListenClient
	)

	ctx, cc = context.WithTimeout(context.Background(), time.Second)
	defer cc()

	if stream, err = testClient.Listen(ctx, &v1.ListenRequest{}); err != nil {
		t.Fatalf("failed to request to listen: %v", err)
	}

	for {
		if _, err = stream.Recv(); err != nil {
			break
		}

		log.Println("data received")

		cc()
	}

	select {
	case <-stream.Context().Done():
	case <-time.After(time.Millisecond * 500):
		t.Fatal("failed to close")
	}

	time.Sleep(time.Millisecond * 100)
}
