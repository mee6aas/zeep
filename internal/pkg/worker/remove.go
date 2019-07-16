package worker

import (
	"context"

	"github.com/pkg/errors"
)

// Remove removes a worker.
func (w *Worker) Remove(ctx context.Context) (err error) {
	if err = w.container.Remove(ctx); err != nil {
		errors.Wrapf(err, "Failed to remove container %s used by the worker %s", w.container.ID(), w.ID())
		return
	}

	if err = w.storage.Remove(); err != nil {
		errors.Wrapf(err, "Failed to remove storage %s used by the worker %s", w.storage.Path(), w.ID())
		return
	}

	return
}
