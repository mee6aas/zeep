package v1_test

import (
	"context"
	"log"
	"net"
	"time"

	v1 "github.com/mee6aas/zeep/pkg/service/invokee/v1"
)

type mockHandle struct {
}

func (h *mockHandle) Connected(ctx context.Context, _ *net.TCPAddr, conn chan<- v1.Task) (err error) {
	log.Println("connected")

	go func() {
		defer close(conn)
		for i := 0; i < 3; i++ {
			select {
			case <-ctx.Done():
				log.Println("canceled")
				return
			default:
				log.Println("send data")
				conn <- v1.Task{}
			}

			time.Sleep(time.Millisecond * 100)
		}
	}()

	return
}

func (h *mockHandle) Disconnected(_ *net.TCPAddr) {
	log.Println("disconnected")
}
