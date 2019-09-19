package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Read gets the worker associated with the specified address.
func Read(addr string) (w *worker.Worker, ok bool) {
	w, ok = workers[addr]

	return
}
