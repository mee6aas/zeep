package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Remove removes a worker from the collection.
func Remove(w *worker.Worker) bool {
	ip := w.IP()
	if _, ok := workers[ip]; !ok {
		return false
	}

	delete(workers, ip)

	return true
}
