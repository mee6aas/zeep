package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy deletes all workers in the pool and stops prewarming.
func (p *Pool) Destroy(ctx context.Context) (err error) {
	ws := p.Entries()
	granted := p.granted

	p.images = []string{}
	p.pendings = map[string]worker.Worker{}
	p.granted = map[string](chan worker.Worker){}

	failed := []string{}

	for _, w := range ws {
		// TODO: go and wait
		if err = w.Remove(ctx); err != nil {
			failed = append(failed, w.ID())
		}
	}

	if err != nil {
		err = errors.Wrapf(err, "Failed to remove %v", failed)
	}

	for _, c := range granted {
		for {
			ok := false

			select {
			case <-c:
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				ok = true
			}

			if ok {
				close(c)
				break
			}
		}
	}

	return
}
