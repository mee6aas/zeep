package pool

import "github.com/mee6aas/zeep/internal/pkg/worker"

// Entries returns all the workers that this pool has.
func (p Pool) Entries() (ws []worker.Worker) {
	for _, e := range p.workers {
		ws = append(ws, append(e[:0:0], e...)...)
	}

	return
}
