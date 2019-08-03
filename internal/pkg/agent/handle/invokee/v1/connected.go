package v1

import (
	"context"
	"errors"
	"net"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

func (h Handle) Connected(
	ctx context.Context,
	addr *net.TCPAddr,
	stream chan<- v1.Task,
) (e error) {
	ok := h.WorkerPool.Grant(addr.IP.String(), &TaskAssigner{
		ctx:    ctx,
		stream: stream,
	}, "1")

	if !ok {
		e = errors.New("Invalid connection")
	}

	return
}
