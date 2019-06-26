package worker

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/container"
	"github.com/mee6aas/zeep/internal/pkg/storage"

	"github.com/pkg/errors"
)

// Remove removes a worker.
func Remove(ctx context.Context, containerID string) (err error) {
	var (
		ok   bool
		conf Config
	)

	if conf, ok = workers[containerID]; !ok {
		// TODO: return NotFound interface
		errors.Errorf("No such worker %s", containerID)
		return
	}

	if err = container.Remove(ctx, containerID); err != nil {
		errors.Wrapf(err, "Failed to remove container %s", containerID)
		return
	}

	delete(workers, containerID)

	if err = storage.Remove(conf.path); err != nil {
		errors.Wrapf(err, "Failed to remove storage %s", conf.path)
		return
	}

	return
}
