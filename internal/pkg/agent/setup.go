package agent

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

// Config holds the configuration for the agent.
type Config struct {
	Acts acts.Config
	Pool pool.Config
}

// Setup initializes agent.
func Setup(ctx context.Context, conf Config) (e error) {
	if isSetup {
		return
	}

	if e = acts.Setup(conf.Acts); e != nil {
		return
	}

	// TODO: this config and opts are testing perposes.
	if workerPool, e = pool.NewPool(ctx, conf.Pool,
		pool.WithEachCPU(0),
		pool.WithEachMem(0),
		pool.WithMaxCPU(0),
		pool.WithMaxMem(0),
	); e != nil {
		return
	}

	isSetup = true

	return
}
