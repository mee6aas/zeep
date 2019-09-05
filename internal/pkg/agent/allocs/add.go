package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Add inserts worker in collection with given username and activity name.
func Add(username string, actName string, w worker.Worker) (ok bool) {
	var (
		es map[string][]worker.Worker
		ws []worker.Worker
	)

	if ok = w.IsAllocated(); !ok {
		return
	}

	if es, ok = allocs[username]; !ok {
		es = make(map[string][]worker.Worker)
	}

	if ws, ok = es[actName]; !ok {
		ws = make([]worker.Worker, 0, 1)
		ok = true
	}

	es[actName] = append(ws, w)
	allocs[username] = es

	return
}
