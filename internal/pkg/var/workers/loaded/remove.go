package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Remove removes a worker from the collection.
func Remove(w *worker.Worker) (ok bool) {
	ip := w.IP()
	if _, ok = workers[ip]; !ok {
		return
	}

	delete(workers, ip)

	return
}
