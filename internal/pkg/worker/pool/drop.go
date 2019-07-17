package pool

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Drop removes worker and restore the resource of pool.
func (p *Pool) Drop(w worker.Worker) (err error) {
	err = w.Remove(context.Background())

	// TODO update used* fields
	// might need to manage fetched workers list

	return
}
