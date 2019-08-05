package worker

import (
	"context"

	"github.com/pkg/errors"
)

// Assign passes task to operator
func (w *Worker) Assign(ctx context.Context, task interface{}) (e error) {
	if !w.IsAllocated() {
		e = errors.New("Task operator not allocated")
		return
	}

	w.isAssigned = true
	e = w.taskAssigner.Assign(ctx, task)

	return
}
