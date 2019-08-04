package agent

import (
	"context"
	"net"
	"os"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

// Config holds the configuration for the agent.
type Config struct {
	AccessPoint string // The IP address that worker connect.

	Acts acts.Config
	Pool pool.Config
}

// Setup initializes agent.
func Setup(ctx context.Context, conf Config) (e error) {
	if isSetup {
		return
	}

	if conf.AccessPoint == "" {
		// do nothing
	} else if host, port, err := net.SplitHostPort(conf.AccessPoint); err == nil {
		os.Setenv("AGENT_HOST", host)
		os.Setenv("AGENT_PORT", port)
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
