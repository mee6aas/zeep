package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Take withdraws worker from collection with specified task id.
func Take(actID string) (w worker.Worker, ok bool) {
	var (
		ws []worker.Worker
	)

	for {
		if ws, ok = allocs[actID]; !ok {
			return
		}

		if len(ws) == 0 {
			ok = false
			return
		}

		w, allocs[actID] = ws[0], ws[1:]

		if len(ws) == 0 {
			delete(allocs, actID)
		}

		if w.IsAllocated() {
			return
		}
	}

}
