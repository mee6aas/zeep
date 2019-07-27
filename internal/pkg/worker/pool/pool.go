package pool

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Pool manages workers
type Pool struct {
	images []string // Images for the worker.
	option Options

	usedCPU float64 // Cpu resources allocated to this pool.
	usedMem uint64  // Amount of memory allocated to this pool in KiB.

	//			image
	workers map[string][]worker.Worker
}

// Config holds the configuration for the pool.
type Config struct {
	Images []string // Images to used by the workers
}

// NewPool creates a new pool based on the given configuration and options.
func NewPool(
	ctx context.Context,
	config Config,
	setters ...Option,
) (
	pool Pool,
	err error,
) {
	args := Options{
		eachCPU: 1,
		eachMem: 1024 * 128,
		maxCPU:  1,
		maxMem:  0,
	}

	for _, setter := range setters {
		setter(&args)
	}

	workers := make(map[string][]worker.Worker)
	for _, image := range config.Images {
		workers[image] = make([]worker.Worker, 0, 32)
	}

	pool = Pool{
		images: config.Images,
		option: args,

		usedCPU: 0,
		usedMem: 0,

		workers: workers,
	}

	for _, image := range config.Images {
		// TODO: go and wait
		if err = pool.alloc(ctx, image); err != nil {
			return
		}
	}

	return
}
