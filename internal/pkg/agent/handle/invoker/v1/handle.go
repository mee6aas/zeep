package v1

import "github.com/mee6aas/zeep/internal/pkg/worker/pool"

// Handle is handler for invoker service
type Handle struct {
	WorkerPool *pool.Pool
}
