package pool

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Drop removes worker and restore the resource of pool.
func (p Pool) Drop(ctx context.Context, w worker.Worker) (e error) {
	e = w.Remove(ctx)

	// TODO update used* fields
	// might need to manage fetched workers list

	return
}
