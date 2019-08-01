package pool

import "github.com/mee6aas/zeep/internal/pkg/worker"

// Entries returns all the workers that this pool has.
func (p Pool) Entries() (ws map[string]worker.Worker) {
	ws = p.pendings

	return
}
