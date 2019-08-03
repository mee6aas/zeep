package agent

import (
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	isSetup bool

	// WorkerPool is set of idle workers
	// TODO: make it configurable
	WorkerPool pool.Pool
)
