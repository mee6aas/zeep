package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Add inserts worker in collection with specified id.
func Add(actID string, w worker.Worker) (ok bool) {
	var (
		ws []worker.Worker
	)

	if ok = w.IsAllocated(); !ok {
		return
	}

	if ws, ok = allocs[actID]; !ok {
		ws = make([]worker.Worker, 0, 1)
		ok = true
	}

	allocs[actID] = append(ws, w)

	return
}
