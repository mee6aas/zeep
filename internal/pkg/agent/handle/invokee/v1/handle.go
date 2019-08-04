package v1

import "github.com/mee6aas/zeep/internal/pkg/worker/pool"

// Handle is handler for invokee service
type Handle struct {
	WorkerPool *pool.Pool
}
