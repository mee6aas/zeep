package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy deletes all workers in the pool and stops prewarming.
func (p *Pool) Destroy(ctx context.Context) (err error) {
	// I think cancel should be first
	//  but that makes removing container fail.
	// So I wait first, but still there is a possibility
	//  that `alloc` called after the wait finished in multithread environment.
	// TODO: resolve it
	//  destoried flag maybe useful
	p.wg.Wait()
	p.cancel()

	ws := p.Entries()
	granted := p.granted

	p.images = make([]string, 0)
	p.pendings = make(map[string]worker.Worker)
	p.granted = make(map[string](chan worker.Worker))

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
