package v1

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

// InvokeRequested is invoked when the invoker requests an activity invoke.
func (h Handle) InvokeRequested(
	ctx context.Context,
	username string,
	actName string,
	actLabel string,
	arg string,
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
		if w, e = h.WorkerPool.Fetch(ctx, a.Runtime); e != nil {
			// TODO: provide details of error
			// possible reasons
			//	- runtime not provided by pool
			//	- lack of resource
			//	- fail to create worker
			e = errors.Wrap(e, "Failed to fetch worker from pool")
			return
		}

		// already checked if the username exists
		actP, _ := acts.PathOf(username)
		// bind activity resources
		w.AddActs(actP)

		loadID, c := assigns.Add()

		// warming worker
		if e = w.Assign(ctx, invokeeV1API.Task{
			Id:   loadID,
			Type: invokeeV1API.TaskType_LOAD,
			Arg:  actName,
		}); e != nil {
			e = errors.Wrap(e, "Failed to assign load task to worker")

			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				if e := w.Remove(ctx); e != nil {
					// TODO: warn error
				}
			}()

			return
		}

		// wait load task finished
		rst := <-c

		// deallocated while invocation
		if rst == nil {
			// TODO:
		}
	}

	// add assignment to list
	invID, c := assigns.Add()

	// assign task to worker
	if e = w.Assign(ctx, invokeeV1API.Task{
		Id:   invID,
		Type: invokeeV1API.TaskType_INVOKE,
		Arg:  arg,
	}); e != nil {
		e = errors.Wrap(e, "Failed to assign task to worker")
		return
	}

	// wait task finished
	rst := <-c

	// deallocated while invocation
	if rst == nil {
		// TODO:
	}

	switch r := rst.(type) {
	case *invokeeV1API.ReportRequest:
		res = &invokerV1.InvokeResponse{
			Result: r.GetResult(),
		}
	default:
		panic(errors.Errorf("Unrecognized report request %v", rst))
	}

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
