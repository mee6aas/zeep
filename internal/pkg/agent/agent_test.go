package agent_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mee6aas/zeep/internal/pkg/agent"
	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/agent/mock/invokee"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	testAddress = "localhost:5121"

	testSetupFailed   = false
	testDestroyFailed = false
)

func TestSetup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Setup(ctx, agent.Config{
		Addr: testAddress,
		Acts: acts.Config{},
		Pool: pool.Config{Images: []string{"mee6aas/runtime-test:latest"}},
	}); err != nil {
		t.Fatalf("Failed to setup agent: %v", err)
		testSetupFailed = true
	}
}

func TestServe(t *testing.T) {
	if testSetupFailed {
		t.Skipf("TestSetup failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	if err := agent.Serve(ctx); (err != nil) &&
		(err != context.DeadlineExceeded) {
		t.Fatalf("Failed to server agent: %v", err)
	}
}

func TestInvokeeListenFail(t *testing.T) {
	if testSetupFailed {
		t.Skipf("TestSetup failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go agent.Serve(ctx)

	i := invokee.Invokee{}
	defer i.Close()

	if err := i.Connect(testAddress); err != nil {
		t.Fatalf("Failed to connect to invokee service: %v", err)
	}

	if err := i.Listen(ctx); err != nil {
		t.Fatalf("Failed to listen from invokee service")
	}

	time.Sleep(time.Millisecond * 100)

	if _, err := i.FetchTask(); err == nil {
		t.Fatalf("Expected to fail to fetch task")
	} else if s, _ := status.FromError(err); s.Code() != codes.PermissionDenied {
		t.Fatalf("Expected to error is PermissionDenied")
	}
}

func TestDestory(t *testing.T) {
	if testSetupFailed {
		t.Skipf("TestSetup failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Destroy(ctx); err != nil {
		testDestroyFailed = true
		t.Fatalf("Failed to destory agent: %v", err)
	}
}
