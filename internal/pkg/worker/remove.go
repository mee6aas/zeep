package worker

import (
	"context"
	"path/filepath"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/storage"

	"github.com/pkg/errors"
)

// Remove removes a worker.
func (w Worker) Remove(ctx context.Context) (e error) {
	if e = w.container.Remove(ctx); e != nil {
		e = errors.Wrapf(e, "Failed to remove container %s used by the worker %s", w.container.ID(), w.ID())
		return
	}

	if e = storage.Unmount(filepath.Join(w.storage.Path(), filepath.Base(api.ActivityResource))); e != nil {
		// not mounted
		if e.Error() != "invalid argument" {
			e = errors.Wrapf(e, "Failed to unmount directory %s used by the worker %s",
				filepath.Join(w.storage.Path(), filepath.Base(api.ActivityResource)), w.ID())
			return
		}
	}

	if e = w.storage.Remove(); e != nil {
		e = errors.Wrapf(e, "Failed to remove storage %s used by the worker %s", w.storage.Path(), w.ID())
		return
	}

	return
}
