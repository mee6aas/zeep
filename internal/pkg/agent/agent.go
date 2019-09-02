package agent

import (
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	isSetup bool

	// workerPool is set of idle workers
	// TODO: make it configurable
	workerPool pool.Pool
)
