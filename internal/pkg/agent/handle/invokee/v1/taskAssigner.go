package v1

import (
	"context"
	"errors"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

type TaskAssigner struct {
	ctx    context.Context
	stream chan<- v1.Task
}

func (ta TaskAssigner) Assign(ctx context.Context, t interface{}) (err error) {
	select {
	case <-ta.ctx.Done():
		err = errors.New("Disconnected")
	case <-ctx.Done():
		err = ctx.Err()
	case ta.stream <- *(t.(*v1.Task)):
	}

	return

}
