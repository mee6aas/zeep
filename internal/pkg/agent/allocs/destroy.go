package allocs

import (
	"context"
	"sync"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Destroy removes allocated workers.
func Destroy(ctx context.Context) (e error) {
	wg := sync.WaitGroup{}

	for _, es := range allocs {
		for _, ws := range es {
			for _, w := range ws {
				wg.Add(1)
				defer func(w worker.Worker) {
					defer wg.Done()
					w.RemoveDetach(ctx)
				}(w)
			}
		}
	}

	wg.Wait()

	return
}
