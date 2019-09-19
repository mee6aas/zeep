package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
	"github.com/pkg/errors"
)

// Add adds a worker into the collection with specified username and activity name as the key.
func Add(username string, actName string, w worker.Worker) (e error) {
	var (
		ok bool
		es map[string][]worker.Worker
		ws []worker.Worker
	)

	if ok = w.IsAllocated(); !ok {
		e = errors.New("not allocated worker")
		return
	}

	if es, ok = workers[username]; !ok {
		es = make(map[string][]worker.Worker)
	}

	if ws, ok = es[actName]; !ok {
		ws = make([]worker.Worker, 0, 1)
	}

	es[actName] = append(ws, w)
	workers[username] = es

	return
}

// TryAdd attempts to add a worker into the collection with specified username and activity name as the key.
// Use the `Add` instead if you are not sure about the preconditions.
func TryAdd(username string, actName string, w worker.Worker) (ok bool) {
	e := Add(username, actName, w)
	ok = e == nil

	return
}
