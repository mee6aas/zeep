package worker_test

import (
	"context"
	"testing"
	"time"

	. "github.com/mee6aas/zeep/internal/pkg/worker"
)

func TestWorkerCreate(t *testing.T) {
	var (
		err    error
		ctx    = context.Background()
		cancel context.CancelFunc
		res    CreateCreatedBody
	)

	ctx, cancel = context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	if res, err = Create(ctx, &Config{
		Image: "notExists",
		Size:  1024,
	}); err == nil {
		t.Fatalf("Expected to fail to create worker")
	}

	_ = res

	// TODO: there are no runtime yet.
}
