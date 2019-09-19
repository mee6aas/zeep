package worker

import (
	"context"

	"github.com/pkg/errors"
)

// TaskAssigner passes task to its task operator.
type TaskAssigner interface {
	Assign(context.Context, interface{}) error
	Close()
}

// Allocate sets the specified task assigner to this worker.
func (w *Worker) Allocate(ta TaskAssigner) (e error) {
	if w.IsAllocated() {
		e = errors.New("Already allocated worker")
		return
	}

	w.isAllocated = true
	w.taskAssigner = ta

	return
}

// Dealloc set isAllocated flag to false
func (w *Worker) Dealloc() { w.isAllocated = false }

// Reallocate sets the specified task assigner to this worker.
func (w *Worker) Reallocate(ta TaskAssigner) (e error) {
	if !w.IsAllocated() {
		e = errors.New("Worker never allocated")
		return
	}

	w.taskAssigner.Close()
	w.taskAssigner = ta

	return
}
