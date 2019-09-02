package pool

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

func (p *Pool) alloc(ctx context.Context, image string) (e error) {
	var (
		ok = false
		w  worker.Worker
	)

	p.allocating.Add(1)
	defer p.allocating.Done()

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

	if e = w.Start(ctx); e != nil {
		w.RemoveDetach(p.ctx)

		return
	}

	p.pendings[w.IP()] = w

	log.WithField("IP", w.IP()).Debug("New worker pended")

	// TODO: update used* fields

	return
}
