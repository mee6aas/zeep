package allocs_test

import (
	"context"
	"testing"

	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

type mockTaskAssigner struct{}

func (m mockTaskAssigner) Assign(context.Context, interface{}) (err error) { return }

func TestAllocs(t *testing.T) {
	var (
		ok bool

		username = "Jerry"
		actName  = "m6s"

		v = "C-137"
		w = worker.Worker{InvokeeVersion: v}
	)

	if ok = allocs.Add(username, actName, w); ok {
		t.Fatal("Expected to fail to add worker")
	}

	if _, ok = allocs.Take(username, actName); ok {
		t.Fatal("Expected to fail to take worker")
	}

	if err := w.Allocate(&mockTaskAssigner{}); err != nil {
		t.Fatal("Failed to allocate a worker")
	}

	if ok = allocs.Add(username, actName, w); !ok {
		t.Fatal("Expected to add worker")
	}

	if w, ok = allocs.Take(username, actName); !ok {
		t.Fatal("Expected to take worker")
	}

	if w.InvokeeVersion != v {
		t.Fatalf("Expected that the taken worker would be the same worker that added but %s", w.InvokeeVersion)
	}
}
