package worker

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/storage"
)

// Remove removes a worker.
func (w Worker) Remove(ctx context.Context) (e error) {
	var (
		errs []string
	)

	if err := w.container.Remove(ctx); err != nil {
		err = errors.Wrapf(err, "Failed to remove container")
		errs = append(errs, err.Error())
	}

	trg := filepath.Join(w.storage.Path(), filepath.Base(api.ActivityResource))

	if err := storage.Unmount(trg); err != nil && err.Error() != "invalid argument" {
		err = errors.Wrapf(err, "Failed to unmount directory %s", trg)
		errs = append(errs, err.Error())
	} else if err := w.storage.Remove(); err != nil {
		err = errors.Wrapf(err, "Failed to remove storage %s", w.storage.Path())
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		e = fmt.Errorf(strings.Join(errs, "\n"))
	}

	return
}

// RemoveDetach removes a worker and warns if failed not return error.
func (w Worker) RemoveDetach(ctx context.Context) {
	if err := w.Remove(ctx); err != nil {
		// warn
	}
}
