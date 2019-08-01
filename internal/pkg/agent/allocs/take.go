package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Take withdraws worker from collection with given username and activity name.
func Take(username string, actName string) (w worker.Worker, ok bool) {
	var (
		es map[string][]worker.Worker
		ws []worker.Worker
	)

	for {
		if es, ok = allocs[username]; !ok {
			return
		}

		if ws, ok = es[actName]; !ok {
			return
		}

		if len(ws) == 0 {
			ok = false
			return
		}

		w, es[actName] = ws[0], ws[1:]

		if len(ws) == 0 {
			delete(es, actName)
		}

		if len(es) == 0 {
			delete(allocs, username)
		}

		allocs[username] = es

		if w.IsAllocated() {
			return
		}
	}

}
