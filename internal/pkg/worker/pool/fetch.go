package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Fetch withdraws a worker in this pool.
func (p *Pool) Fetch(ctx context.Context, image string) (w worker.Worker, e error) {
	var (
		ok = false
	)

	for _, img := range p.images {
		ok = img == image
	}
	if !ok {
		e = errors.New("not found")
		return
	}

	select {
	case w = <-p.granted[image]:
	case <-ctx.Done():
	}

	if e = ctx.Err(); e != nil {
		return
	}

	if ok = w.IsAllocated(); !ok {
		e = errors.New("Fail")
		return
	}

	go p.alloc(p.ctx, image)

	return
}
