package agent

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
	"github.com/pkg/errors"
)

// Destroy removes agent resources.
func Destroy(ctx context.Context) (e error) {
	if !isSetup {
		return
	}

	// TODO: currently only last error is preserved.

	if err := acts.Destroy(ctx); err != nil {
		e = errors.Wrap(err, "Failed to destory acts")
	}

	if err := allocs.Destroy(ctx); err != nil {
		e = errors.Wrap(err, "Failed to destory allocs")
	}

	if err := workerPool.Destroy(ctx); err != nil {
		e = errors.Wrap(err, "Failed to destory pool")
	}

	isSetup = false

	return
}
