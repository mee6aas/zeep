package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Add adds a worker into the collection.
// The IP of the worker used as a key.
func Add(w *worker.Worker) bool {
	ip := w.IP()
	if _, ok := workers[ip]; ok {
		return false
	}

	workers[ip] = w

	return true
}
