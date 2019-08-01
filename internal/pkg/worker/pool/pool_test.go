package pool_test

import (
	"context"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/worker"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

type mockTA struct{}

func (ta mockTA) Assign(context.Context, interface{}) (err error) { return }

var (
	testImage = "golang:1.12"
	testPool  pool.Pool

	testNewPoolFailed = false
)

func TestNewPoolFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := pool.NewPool(ctx, pool.Config{
		Images: []string{"not exist"},
	}); err == nil {
		t.Fatal("Expected to fail to create a new pool.")
	}
}

func TestNewPool(t *testing.T) {
	var (
		err error
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if testPool, err = pool.NewPool(ctx, pool.Config{
		Images: []string{testImage},
	}); err != nil {
		testNewPoolFailed = true
		t.Fatalf("Failed to create pool: %v", err)
	}
}

func TestFetch(t *testing.T) {
	if testNewPoolFailed {
		t.Skip("TestNewPool failed")
	}

	var (
		err error
		w   worker.Worker
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for _, w := range testPool.Entries() {
		testPool.Grant(w.Container().IP(), &mockTA{}, "C-137")
	}

	if w, err = testPool.Fetch(ctx, testImage); err != nil {
		t.Fatalf("Failed to fetch worker from pool: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		w.Remove(ctx)
	}()

	if !w.Container().IsExists(ctx) {
		t.Fatal("Expected that the container exists.")
	}

	t.Logf("Worker %s fetched", w.ID())
}

func TestPrewarming(t *testing.T) {
	var (
		id = ""
	)

	if testNewPoolFailed {
		t.Skip("TestNewPool failed")
	}

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 200)

		for _, w := range testPool.Entries() {
			id = w.ID()
			break
		}

		if id != "" {
			break
		}
	}

	if id != "" {
		t.Logf("Worker %s prewarmed", id)
	} else {
		t.Fatal("Expected that the worker had created.")
	}
}

func TestDestory(t *testing.T) {
	if testNewPoolFailed {
		t.Skip("TestNewPool failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := testPool.Destroy(ctx); err != nil {
		t.Fatalf("Failed to destroy the pool: %v", err)
	}
}
