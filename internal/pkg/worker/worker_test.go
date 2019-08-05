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
		err error
		w   Worker
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	if w, err = NewWorker(ctx, worker.Config{
		Image: "notExists",
		Size:  1024,
	}); err == nil {
		t.Fatalf("Expected to fail to create worker")
	}

	_ = w

	// TODO: there are no runtime yet.
}
