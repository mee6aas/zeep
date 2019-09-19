package workers

import (
	"context"
	"sync"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy removes all workers in the collection.
func Destroy(ctx context.Context) (e error) {
	wg := sync.WaitGroup{}

	for _, es := range workers {
		for _, ws := range es {
			for _, w := range ws {
				wg.Add(1)
				go func(w worker.Worker) {
					defer wg.Done()
					w.RemoveDetach(ctx)
				}(w)
			}
		}
	}

	wg.Wait()

	return
}
