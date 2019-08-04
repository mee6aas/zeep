package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

func (p *Pool) alloc(ctx context.Context, image string) (e error) {
	var (
		ok = false
		w  worker.Worker
	)

	p.wg.Add(1)
	defer p.wg.Done()

	for _, img := range p.images {
		ok = img == image
	}
	if !ok {
		e = errors.New("Image not provided by pool")
		return
	}

	if w, e = worker.NewWorker(ctx, worker.Config{
		Image: image,
		// TODO: isolation
		// Size:
	}); e != nil {
		return
	}

	c := w.Container()
	if e = c.Start(ctx); e != nil {
		// TODO: handle orphan worker
		_ = w.Remove(ctx)

		return
	}

	p.pendings[c.IP()] = w

	// TODO: update used* fields

	return
}
