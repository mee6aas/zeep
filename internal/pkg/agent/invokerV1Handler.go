package agent

import (
	"context"
	"time"

	"github.com/pkg/errors"

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
	username string,
	actName string,
	actLabel string,
) (
	res *invokerV1.InvokeResponse,
	e error,
) {
	var (
		ok bool
		a  activity.Activity
		w  worker.Worker
	)

	// read activity manifest
	if a, ok = acts.Read(username, actName); !ok {
		e = errors.New("Not found")
		return
	}

	// get worker
	if w, ok = allocs.Take(username, actName); !ok {
		// warm worker not exists
		if w, e = workerPool.Fetch(a.Runtime); e != nil {
			// TODO: provide details of error
			// possible reasons
			//	- runtime not provided by pool
			//	- lack of resource
			//	- fail to create worker
			e = errors.Wrap(e, "Failed to fetch worker from pool")
			return
		}

		// TODO: wait for worker allocated

		// warming worker
		// TODO: move activity resource to container
		// w.Load(username + "/" + actName)?
		if e = w.Assign(ctx, invokeeV1API.Task{
			Id:   username + "/" + actName,
			Type: invokeeV1API.TaskType_LOAD,
		}); e != nil {
			e = errors.Wrap(e, "Worker failed to load activity")

			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				if e := w.Remove(ctx); e != nil {
					// TODO: warn error
				}
			}()

			return
		}
	}

	// add assignment to list
	invID, c := assigns.Add()

	// assign task to worker
	if e = w.Assign(ctx, invokeeV1API.Task{
		Id:   invID,
		Type: invokeeV1API.TaskType_INVOKE,
	}); e != nil {
		e = errors.Wrap(e, "Failed to assign task to worker")
		return
	}

	// wait task finished
	rst := <-c

	// deallocated while invocation
	if rst == nil {
	}

	// TODO: implement version converter
	// but it is ok currently because there is only v1.
	res, _ = rst.(*invokerV1.InvokeResponse)

	// resolve worker
	w.Resolve()

	// maintain worker
	if ok = allocs.Add(username, actName, w); !ok {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			if e := w.Remove(ctx); e != nil {
				// TODO: warn error
			}
		}()
	}

	return
}

func (h invokerV1Handle) RegisterRequested(
	_ context.Context,
	username string,
	actDirPath string,
) (e error) {
	var (
		a activity.Activity
	)

	a, e = activity.UnmarshalFromDir(actDirPath)
	e = acts.Add(username, a.Name, actDirPath)

	return
}
