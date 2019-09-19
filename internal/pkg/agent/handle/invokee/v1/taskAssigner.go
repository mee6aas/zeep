package v1

import (
	"context"
	"errors"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// TaskAssigner holds stream for assigning a task to the worker.
type TaskAssigner struct {
	ctx    context.Context
	stream chan<- v1.Task
}

// Assign sends the specified task to the worker.
func (ta TaskAssigner) Assign(ctx context.Context, t interface{}) (e error) {
	select {
	case <-ta.ctx.Done():
		e = errors.New("Disconnected")
	case <-ctx.Done():
		e = ctx.Err()
	case ta.stream <- (t.(v1.Task)):
	}

	return

}

// Close closes the connected channel
func (ta TaskAssigner) Close() {
	close(ta.stream)
}
