package v1

import (
	"context"
	"errors"
	"net"

	"github.com/mee6aas/zeep/internal/pkg/agent/assistant/handover"
	workers "github.com/mee6aas/zeep/internal/pkg/var/workers/loaded"
	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Connected is invoked when invokee client connected.
func (h Handle) Connected(
	ctx context.Context,
	addr *net.TCPAddr,
	stream chan<- v1.Task,
) (e error) {
	ip := addr.IP.String()

	// initial connection
	if ok := h.WorkerPool.Grant(ip, &TaskAssigner{
		ctx:    ctx,
		stream: stream,
	}, "1"); ok {
		return
	}

	// handover
	if ok := handover.IsReserved(addr.IP.String()); !ok {
		e = errors.New("Invalid connection")
		return
	}

	w, ok := workers.Read(ip)
	if !ok {
		e = errors.New("Invalid handover procedure")
		return
	}

	w.Reallocate(&TaskAssigner{
		ctx:    ctx,
		stream: stream,
	})

	_ = handover.Resolve(ip)

	return
}
