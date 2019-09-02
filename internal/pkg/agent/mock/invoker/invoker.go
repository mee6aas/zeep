package invoker

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// Invoker is mock invoker.
type Invoker struct {
	Conn   *grpc.ClientConn
	Client invokerV1API.InvokerClient
}

// Connect connects to agent invoker service.
func (i *Invoker) Connect(address string) (e error) {
	if i.Conn, e = grpc.Dial(address, grpc.WithInsecure()); e != nil {
		return
	}

	i.Client = invokerV1API.NewInvokerClient(i.Conn)

	return
}

// Close closes connection.
func (i *Invoker) Close() (e error) {
	e = i.Conn.Close()

	return
}

// Add requests add service.
func (i *Invoker) Add(
	ctx context.Context,
	username string,
	actName string,
	actDirPath string,
) (e error) {
	_, e = i.Client.Add(ctx, &invokerV1API.AddRequest{
		Username: username,
		ActName:  actName,
		Path:     actDirPath,
		Method:   invokerV1API.AddMethod_LOCAL,
	})

	return
}

// Invoke requests invoke service.
func (i *Invoker) Invoke(
	ctx context.Context,
	username string,
	actName string,
	arg string,
) (rst string, e error) {
	var (
		res *invokerV1API.InvokeResponse
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if res, e = i.Client.Invoke(ctx, &invokerV1API.InvokeRequest{
			Username: username,
			ActName:  actName,
			Arg:      arg,
		}); e != nil {
			return
		}
		rst = res.GetResult()
	}()

	wg.Wait()

	return
}
