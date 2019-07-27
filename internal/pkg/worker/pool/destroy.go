package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy deletes all workers in the pool and stops prewarming.
func (p *Pool) Destroy(ctx context.Context) (err error) {
	ws := p.Entries()

	p.images = []string{}
	p.workers = map[string][]worker.Worker{}

	failed := []string{}

	for _, w := range ws {
		if err = w.Remove(ctx); err != nil {
			failed = append(failed, w.ID())
		}
	}

	if err != nil {
		err = errors.Wrapf(err, "Failed to remove %v", failed)
	}

	return
}
