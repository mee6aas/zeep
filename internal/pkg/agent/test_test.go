package agent

import "github.com/mee6aas/zeep/internal/pkg/worker/pool"

func WorkerPool() (wp *pool.Pool) {
	wp = &workerPool
	return
}
