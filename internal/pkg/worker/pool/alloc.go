package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

func (p Pool) alloc(ctx context.Context, image string) (err error) {
	var (
		ok bool
		w  worker.Worker
		ws []worker.Worker
	)

	if ws, ok = p.workers[image]; !ok {
		err = errors.New("not found")
		return
	}

	if w, err = worker.NewWorker(ctx, worker.Config{
		Image: image,
		// TODO: isolation
		// Size:
	}); err != nil {
		return err
	}

	cont := w.Container()
	if err = cont.Start(ctx); err != nil {
		// TODO: handle orphan worker
		_ = w.Remove(ctx)

		return err
	}

	ws = append(ws, w)

	// TODO: update used* fields

	return
}
