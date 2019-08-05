package pool

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy deletes all workers in the pool and stops prewarming.
func (p *Pool) Destroy(ctx context.Context) (e error) {
	// I think cancel should be first
	//  but that makes removing container fail.
	// So I wait first, but still there is a possibility
	//  that `alloc` called after the wait finished in multithread environment.
	// TODO: resolve it
	//  destoried flag maybe useful
	p.allocating.Wait()
	p.cancel()

	var (
		ws = p.pendings
		gs = p.granted
	)

	p.images = make([]string, 0)
	p.pendings = make(map[string]worker.Worker)
	p.granted = make(map[string](chan worker.Worker))

	// remove pendings and grantedworkers
	failed := make([]string, 0, len(ws))
	{
		var (
			wg  = sync.WaitGroup{}
			mtx = sync.Mutex{}
		)

		for _, w := range ws {
			wg.Add(1)
			go func(w worker.Worker) {
				defer wg.Done()
				if err := w.Remove(ctx); err != nil {
					mtx.Lock()
					e = err
					failed = append(failed, w.ID())
					mtx.Unlock()
				}
			}(w)
		}

		wg.Wait()
	}

	if e != nil {
		e = errors.Wrapf(e, "Failed to remove %v", failed)
	}

	// clear channels
	// GC would not collect channels since granted workers are in the
	//  each goroutine and holds channels.
	for _, c := range gs {
		for {
			ok := false

			select {
			case <-c:
			case <-ctx.Done():
				if ctx.Err() != nil {
					e = errors.Wrapf(e, ctx.Err().Error())
				}
				e = errors.Wrapf(e, "Context canceled during destructing agent")
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
