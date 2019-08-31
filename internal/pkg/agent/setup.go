package agent

import (
	"context"
	"net"
	"os"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

// Config holds the configuration for the agent.
type Config struct {
	Addr string // The address agent to serve
	Acts acts.Config
	Pool pool.Config
}

// Setup initializes agent.
func Setup(ctx context.Context, conf Config) (e error) {
	if isSetup {
		return
	}

	addr = conf.Addr

	if conf.Addr == "" {
		// do nothing
	} else if host, port, err := net.SplitHostPort(conf.Addr); err == nil {
		os.Setenv(api.AgentHostEnvKey, host)
		os.Setenv(api.AgentPortEnvKey, port)
	} else {
		e = err
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
