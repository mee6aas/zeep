package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/agent"
	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	testSetupFailed = false
)

func TestSetup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Setup(ctx, agent.Config{
		Acts: acts.Config{},
		Pool: pool.Config{Images: []string{"golang:1.12"}},
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

	if err := agent.Serve(ctx, "localhost:5122"); (err != nil) &&
		(err != context.DeadlineExceeded) {
		t.Fatalf("Failed to server agent: %v", err)
	}
}

func TestDestory(t *testing.T) {
	if testSetupFailed {
		t.Skipf("TestSetup failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Destroy(ctx); err != nil {
		t.Fatalf("Failed to destory agent: %v", err)
	}
}
