package invokee

import (
	"context"

	"google.golang.org/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Invokee is mock invokee.
type Invokee struct {
	Conn   *grpc.ClientConn
	Client invokeeV1API.InvokeeClient

	TaskStream    invokeeV1API.Invokee_ListenClient
	TaskStreamErr error
	TaskChan      chan *invokeeV1API.Task

	ListenCtx    context.Context
	ListenCancel context.CancelFunc
}

// Connect connects to agent invokee service.
func (i *Invokee) Connect(address string) (e error) {
	if i.Conn, e = grpc.Dial(address, grpc.WithInsecure()); e != nil {
		return
	}

	i.Client = invokeeV1API.NewInvokeeClient(i.Conn)

	return
}

// Close closes connection.
func (i *Invokee) Close() (e error) {
	i.ListenCancel()
	close(i.TaskChan)
	e = i.Conn.Close()

	return
}

// Listen requests listen service.
func (i *Invokee) Listen(ctx context.Context) (e error) {
	i.ListenCtx, i.ListenCancel = context.WithCancel(ctx)

	if i.TaskStream, e = i.Client.Listen(i.ListenCtx, &invokeeV1API.ListenRequest{}); e != nil {
		return
	}

	i.TaskChan = make(chan *invokeeV1API.Task)

	go func() {
		for {
			t, e := i.TaskStream.Recv()
			if e != nil {
				i.TaskStreamErr = e
				return
			}
			i.TaskChan <- t
		}
	}()

	return
}

// FetchTask returns task if exists.
func (i *Invokee) FetchTask() (t *invokeeV1API.Task, e error) {
	select {
	case t = <-i.TaskChan:
	default:
		e = i.TaskStreamErr
	}

	return
}

// Report send report request to invokee service.
func (i *Invokee) Report(
	ctx context.Context,
	id string,
	rst string,
	isErr bool,
) (e error) {
	_, e = i.Client.Report(ctx,
		&invokeeV1API.ReportRequest{
			Id:      id,
			Result:  rst,
			IsError: isErr,
		},
	)

	return
}
