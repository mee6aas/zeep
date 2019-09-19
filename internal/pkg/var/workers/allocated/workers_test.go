package workers_test

import (
	"context"
	"testing"

	workers "github.com/mee6aas/zeep/internal/pkg/var/workers/allocated"
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

type mockTaskAssigner struct{}

func (m mockTaskAssigner) Assign(context.Context, interface{}) (err error) { return }
func (m mockTaskAssigner) Close()                                          {}

func TestAllocs(t *testing.T) {
	var (
		username = "Jerry"
		actName  = "m6s"

		v = "C-137"
		w = worker.Worker{InvokeeVersion: v}
	)

	if err := workers.Add(username, actName, w); err == nil {
		t.Fatal("Expected to fail to add worker")
	}

	if _, err := workers.Take(username, actName); err == nil {
		t.Fatal("Expected to fail to take worker")
	}

	if err := w.Allocate(&mockTaskAssigner{}); err != nil {
		t.Fatal("Failed to allocate a worker")
	}

	if err := workers.Add(username, actName, w); err != nil {
		t.Fatal("Expected to insert a worker")
	}

	w, err := workers.Take(username, actName)
	if err != nil {
		t.Fatal("Expected to take a worker")
	}

	if w.InvokeeVersion != v {
		t.Fatalf("Expected that the taken worker would be the same worker that added but %s", w.InvokeeVersion)
	}
}
