package worker

import "context"

// Start starts container used by this container use.
func (w *Worker) Start(ctx context.Context) (e error) {
	e = w.container.Start(ctx)
	return
}
