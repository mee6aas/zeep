package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
	"github.com/pkg/errors"
)

// Take withdraws the worker associated with the specified username and activity name from the collection.
func Take(username string, actName string) (w worker.Worker, e error) {
	var (
		ok bool
		es map[string][]worker.Worker
		ws []worker.Worker
	)

	for {
		if es, ok = workers[username]; !ok {
			e = errors.New("username not exists")
			return
		}

		if ws, ok = es[actName]; !ok {
			e = errors.New("activity name not exists")
			return
		}

		if len(ws) == 0 {
			e = errors.New("empty collection")
			return
		}

		w, es[actName] = ws[0], ws[1:]

		if len(ws) == 0 {
			delete(es, actName)
		}

		if len(es) == 0 {
			delete(workers, username)
		}

		workers[username] = es

		// deallocated while in the collection
		if w.IsAllocated() {
			return
		}
	}

}

// TryTake attempts to withdraw a worker from the collection with the specified username and activity name.
// Use the `Take` instead if you are not sure about the preconditions.
func TryTake(username string, actName string) (w worker.Worker, ok bool) {
	w, e := Take(username, actName)
	ok = e == nil

	return
}
