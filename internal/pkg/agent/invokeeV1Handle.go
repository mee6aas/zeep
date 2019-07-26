package agent

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mee6aas/zeep/internal/pkg/agent/assigns"
	v1 "github.com/mee6aas/zeep/pkg/service/invokee/v1"
)

type invokeeV1TaskAssigner struct {
	ctx    context.Context
	stream chan<- v1.Task
}

func (ta invokeeV1TaskAssigner) Assign(ctx context.Context, t interface{}) (err error) {
	select {
	case <-ta.ctx.Done():
		err = errors.New("Disconnected")
	case <-ctx.Done():
		err = ctx.Err()
	case ta.stream <- *(t.(*v1.Task)):
	}

	return

}

type invokeeV1Handle struct{}

func (h invokeeV1Handle) Connected(
	ctx context.Context,
	addr *net.TCPAddr,
	stream chan<- v1.Task,
) (err error) {
	for _, w := range workerPool.Entries() {
		cont := w.Container()
		if cont.IP() == addr.IP.String() {
			w.InvokeeVersion = "1"
			w.Allocate(&invokeeV1TaskAssigner{
				ctx:    ctx,
				stream: stream,
			})

			return
		}
	}

	return errors.New("Invalid connection")
}

func (h invokeeV1Handle) Disconnected(_ *net.TCPAddr) {

}

func (h invokeeV1Handle) Reported(req *v1.ReportRequest) (err error) {
	id := req.GetId()

	if ok := assigns.Report(id, req); !ok {
		err = status.Error(codes.NotFound, "Invocation ID not found")
		return
	}

	return
}
