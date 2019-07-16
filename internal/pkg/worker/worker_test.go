package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/worker"
	. "github.com/mee6aas/zeep/internal/pkg/worker"
)

func TestWorkerCreate(t *testing.T) {
	var (
		err    error
		ctx    = context.Background()
		cancel context.CancelFunc
		wkr    Worker
	)

	ctx, cancel = context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	if wkr, err = NewWorker(ctx, worker.Config{
		Image: "notExists",
		Size:  1024,
	}); err == nil {
		t.Fatalf("Expected to fail to create worker")
	}

	_ = wkr

	// TODO: there are no runtime yet.
}
