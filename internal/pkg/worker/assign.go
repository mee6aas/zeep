package worker

import (
	"context"

	"github.com/pkg/errors"
)

// Assign passes task to operator
func (w *Worker) Assign(ctx context.Context, task interface{}) (err error) {
	if !w.IsAllocated() {
		err = errors.New("Task operator not allocated")
		return
	}

	w.isAssigned = true
	err = w.taskAssigner.Assign(ctx, task)

	if err != nil {
		w.Remove(ctx)
	}

	return
}
