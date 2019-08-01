package v1_test

import (
	"context"
	"log"
	"os"
	"testing"

	"google.golang.org/grpc"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	grpcServer "github.com/mee6aas/zeep/pkg/protocol/grpc"
	inokeeV1 "github.com/mee6aas/zeep/pkg/service/invokee/v1"
)

const (
	testServerAddress = "localhost:5122"
)

var (
	testGrpcServer   = grpc.NewServer()
	testServerHandle = &mockHandle{}
	ctxTestServer    context.Context
	stopTestServer   context.CancelFunc

	testConn   *grpc.ClientConn
	testClient v1.InvokeeClient
)

func TestMain(m *testing.M) {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctxTestServer, stopTestServer = context.WithCancel(context.Background())
	v1.RegisterInvokeeServer(testGrpcServer, inokeeV1.NewInvokeeAPIServer(testServerHandle))

	go grpcServer.Serve(ctxTestServer, testGrpcServer, testServerAddress)

	if testConn, err = grpc.Dial(testServerAddress, grpc.WithInsecure(), grpc.WithWaitForHandshake()); err != nil {
		panic(err)
	}

	testClient = v1.NewInvokeeClient(testConn)

	log.Println("Inokee API client created")

	code := m.Run()
	stopTestServer()
	os.Exit(code)
}
