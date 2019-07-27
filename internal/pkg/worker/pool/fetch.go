package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Fetch fetches a worker in this pool.
func (p *Pool) Fetch(image string) (w worker.Worker, err error) {
	var (
		ok bool
		ws []worker.Worker
	)

	if ws, ok = p.workers[image]; !ok {
		err = errors.New("not found")
		return
	}

	if len(ws) == 0 {
		if err = p.alloc(context.Background(), image); err != nil {
			return
		}
	}

	w, p.workers[image] = ws[0], ws[1:]

	go p.alloc(context.Background(), image)

	return
}
