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

// Register requests register service.
func (i *Invoker) Register(
	ctx context.Context,
	username string,
	actDirPath string,
) (e error) {
	_, e = i.Client.Register(ctx, &invokerV1API.RegisterRequest{
		Username: username,
		Path:     actDirPath,
		Method:   invokerV1API.RegisterMethod_LOCAL,
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
			Target:   &invokerV1API.Activity{Name: actName},
			Arg:      arg,
		}); e != nil {
			return
		}
		rst = res.GetResult()
	}()

	wg.Wait()

	return
}
