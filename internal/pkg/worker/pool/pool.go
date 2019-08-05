package pool

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Pool manages workers
type Pool struct {
	images []string // Images for the worker.
	option Options

	usedCPU float64 // Cpu resources allocated to this pool.
	usedMem uint64  // Amount of memory allocated to this pool in KiB.

	// created but not allocated workers
	//             IP
	pendings map[string]worker.Worker

	//          image
	granted map[string](chan worker.Worker)

	// holds the workers under allocating.
	allocating *sync.WaitGroup

	// life of pool
	ctx    context.Context
	cancel context.CancelFunc
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
	p Pool,
	e error,
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

	granted := make(map[string](chan worker.Worker))
	for _, image := range config.Images {
		granted[image] = make(chan worker.Worker)
	}

	p = Pool{
		images: config.Images,
		option: args,

		usedCPU: 0,
		usedMem: 0,

		pendings: make(map[string]worker.Worker),
		granted:  granted,

		allocating: &sync.WaitGroup{},
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())

	{
		wg := sync.WaitGroup{}
		for _, img := range config.Images {
			wg.Add(1)
			go func(img string) {
				defer wg.Done()
				if e = p.alloc(ctx, img); e != nil {
					p.cancel()
				}
			}(img)
		}
		wg.Wait()
	}

	if e != nil {
		e = errors.Wrap(e, "Failed while allocating")
	}

	return
}
