package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Take withdraws worker from collection with specified id.
func Take(id string) (w worker.Worker, ok bool) {
	var (
		ws []worker.Worker
	)

	for {
		if ws, ok = allocs[id]; !ok {
			return
		}

		if len(ws) == 0 {
			ok = false
			return
		}

		w, allocs[id] = ws[0], ws[1:]

		if len(ws) == 0 {
			delete(allocs, id)
		}

		if w.IsAllocated() {
			return
		}
	}

}
