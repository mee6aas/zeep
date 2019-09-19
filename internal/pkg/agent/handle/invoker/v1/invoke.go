package v1

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"

	acts "github.com/mee6aas/zeep/internal/pkg/var/activities"
	assigns "github.com/mee6aas/zeep/internal/pkg/var/assignments"
	allocatedWorkers "github.com/mee6aas/zeep/internal/pkg/var/workers/allocated"
	invokedWorkers "github.com/mee6aas/zeep/internal/pkg/var/workers/invoked"
	loadedWorkers "github.com/mee6aas/zeep/internal/pkg/var/workers/loaded"
	"github.com/mee6aas/zeep/internal/pkg/worker"
	"github.com/mee6aas/zeep/pkg/activity"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"
)

// InvokeRequested is invoked when the invoker requests an activity invoke.
func (h Handle) InvokeRequested(
	ctx context.Context,
	addr *net.TCPAddr,
	username string,
	actName string,
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

	if actName == "" {
		e = errors.New("Empty activity name")
		return
	}

	// the username is omitted for the invoked worker which invokes another
	if username == "" {
		if username, ok = assigns.GetAssigneeFromIP(addr.IP.String()); !ok {
			e = errors.New("Username not found")
			return
		}
	}

	// read activity manifest
	if a, ok = acts.Read(username, actName); !ok {
		e = errors.New("Not found")
		return
	}

	// get worker
	if w, ok = allocatedWorkers.TryTake(username, actName); !ok {
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

		loadID, c := assigns.Add(w.IP(), username)

		if ok = loadedWorkers.Add(&w); !ok {
			// never happens
			e = errors.New("Already loaded worker")
			return
		}

		// warming worker
		if e = w.Assign(ctx, invokeeV1API.Task{
			Id:   loadID,
			Type: invokeeV1API.TaskType_LOAD,
			Arg:  actName,
		}); e != nil {
			e = errors.Wrap(e, "Failed to assign task for load to worker")

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

		if ok = loadedWorkers.Remove(&w); !ok {
			// never happens
			e = errors.New("Worker never loaded")
			return
		}

		// deallocated while invocation
		if rst == nil {
			// TODO:
		}
	}

	// invoke
	{
		ip := w.IP()
		if ip == "" {
			e = errors.New("IP of the worker is not decided")
			return
		}

		// add assignment to list
		invID, c := assigns.Add(ip, username)

		if ok = invokedWorkers.Add(&w); !ok {
			// never happens
			e = errors.New("Already invoked worker")
			return
		}

		// assign a task to the worker
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

		if ok = invokedWorkers.Remove(&w); !ok {
			// never happens
			e = errors.New("Worker never invoked")
			return
		}

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
	}

	// resolve a worker
	w.Resolve()

	// maintain a worker
	if ok = allocatedWorkers.TryAdd(username, actName, w); !ok {
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
