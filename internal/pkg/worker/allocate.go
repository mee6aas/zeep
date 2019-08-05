package worker

import (
	"context"

	"github.com/pkg/errors"
)

// TaskAssigner passes task to its task operator.
type TaskAssigner interface {
	Assign(context.Context, interface{}) error
}

// Allocate set task assigner to this worker.
func (w *Worker) Allocate(ta TaskAssigner) (e error) {
	if w.IsAllocated() {
		e = errors.New("Already allocated worker")
		return
	}

	w.isAllocated = true
	w.taskAssigner = ta

	return
}
