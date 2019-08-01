package agent

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
	"github.com/mee6aas/zeep/internal/pkg/agent/assigns"
	"github.com/mee6aas/zeep/internal/pkg/worker"
	"github.com/mee6aas/zeep/pkg/activity"
	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

type invokerV1Handle struct{}

func (h invokerV1Handle) InvokeRequested(
	ctx context.Context,
	req *invokerV1.InvokeRequest,
) (
	res *invokerV1.InvokeResponse,
	err error,
) {
	var (
		ok bool
		a  activity.Activity
		w  worker.Worker
	)

	trg := req.GetTarget()

	for _, act := range acts.EntriesInUsername(req.GetUsername()) {
		if act.Name == trg.GetName() {
			a = act
			break
		}
	}

	if a.ID == "" {
		err = errors.New("Not found")
		return
	}

	invID, c := assigns.Add()

	for {
		if w, ok = allocs.Take(a.ID); !ok {
			break
		}

		if err = w.Assign(ctx, invokeeV1API.Task{
			Id:   invID,
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

func (h invokerV1Handle) RegisterRequested(
	_ context.Context,
	username string,
	actDirPath string,
) (e error) {
	var (
		id uuid.UUID
	)

	if id, e = uuid.NewRandom(); e != nil {
		return
	}

	e = acts.Add(username, id.String(), actDirPath)

	return
}
