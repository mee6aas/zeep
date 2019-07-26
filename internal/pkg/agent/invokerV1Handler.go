package agent

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
	"github.com/mee6aas/zeep/internal/pkg/agent/assigns"
	"github.com/mee6aas/zeep/internal/pkg/worker"
	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type invokerV1Handle struct{}

func (h invokerV1Handle) Requested(
	ctx context.Context,
	req *invokerV1.InvokeRequest,
) (
	res *invokerV1.InvokeResponse,
	err error,
) {
	var (
		ok bool
		w  worker.Worker
	)

	trg := req.GetTarget()
	id, c := assigns.Add()

	for {
		if w, ok = allocs.Take(trg.GetId()); !ok {
			break
		}

		if err = w.Assign(ctx, invokeeV1API.Task{
			Id:   id,
			Type: invokeeV1API.TaskType_INVOKE,
		}); err != nil {
			break
		}

		rst := <-c

		if rst == nil {
			// deallocated while invocation
		}

		// TODO: implement version converter
		// but it is ok currently because there is only v1.
		res, _ = rst.(*invokerV1.InvokeResponse)

		w.Resolve()

		return
	}

	return
}
