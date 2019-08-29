// +build int

package agent_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/agent"
	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"

	mockInvokee "github.com/mee6aas/zeep/internal/pkg/agent/mock/invokee"
	mockInvoker "github.com/mee6aas/zeep/internal/pkg/agent/mock/invoker"
)

func TestIntegrationWithMock(t *testing.T) {
	const (
		testAddr        = "172.17.0.1:5210" // docker default ip
		testUsername    = "Jerry"
		testActName     = "empty"
		testActDirPath  = "./testdata/empty"
		testActArg      = "Two strokes off"
		testExpectedRst = "kill Jerry"
	)

	var (
		testRst string
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Setup(ctx, agent.Config{
		Acts: acts.Config{},
		Pool: pool.Config{Images: []string{"mee6aas/runtime-test:latest"}},
	}); err != nil {
		t.Fatalf("Failed to setup agent: %v", err)
	}

	go agent.Serve(ctx, testAddress)
	defer func() {
		if err := agent.Destroy(ctx); err != nil {
			t.Logf("Failed to destory agent: %v", err)
		}
	}()

	for k, _ := range agent.WorkerPool().Entries() {
		agent.WorkerPool().ChangePendingKey("127.0.0.1", k)
		break
	}

	invkee := mockInvokee.Invokee{}
	defer invkee.Close()

	if err := invkee.Connect(testAddress); err != nil {
		t.Fatalf("Failed to connect to invokee service: %v", err)
	}

	if err := invkee.Listen(ctx); err != nil {
		t.Fatalf("Failed to listen from invokee service")
	}

	time.Sleep(time.Millisecond * 100)

	if _, err := invkee.FetchTask(); err != nil {
		t.Fatalf("Failed to fetch task.")
	}

	invker := mockInvoker.Invoker{}
	defer invker.Close()

	if err := invker.Connect(testAddress); err != nil {
		t.Fatalf("Failed to connect to invoker service: %v", err)
	}

	if err := invker.Register(ctx, testUsername, testActName, testActDirPath); err != nil {
		t.Fatalf("Failed to request for register an activity to agent: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if rst, err := invker.Invoke(ctx, testUsername, testActName, testActArg); err != nil {
			t.Fatalf("Failed to invoke an activity: %v", err)
		} else {
			testRst = rst
		}
	}()

	time.Sleep(time.Millisecond * 100)

	if task, err := invkee.FetchTask(); err != nil {
		t.Fatalf("Failed to fetch task: %v", err)
	} else if task == nil {
		t.Fatalf("Expected that the LOAD task is exists.")
	} else if task.GetType() != invokeeV1API.TaskType_LOAD {
		t.Fatalf("Expected that the task type is LOAD")
	} else if err := invkee.Report(ctx, task.GetId(), "", false); err != nil {
		t.Fatalf("Failed to report for LOAD task: %v", err)
	}

	time.Sleep(time.Millisecond * 100)

	if task, err := invkee.FetchTask(); err != nil {
		t.Fatalf("Failed to fetch task: %v", err)
	} else if task == nil {
		t.Fatalf("Expected that the INVOKE task is exists.")
	} else if task.GetType() != invokeeV1API.TaskType_INVOKE {
		t.Fatalf("Expected that the task type is INVOKE")
	} else if err := invkee.Report(ctx, task.GetId(), testExpectedRst, false); err != nil {
		t.Fatalf("Failed to report for LOAD task: %v", err)
	}

	wg.Wait()

	if testRst != testExpectedRst {
		t.Fatalf("Expected that the result is %s but is %s",
			testExpectedRst,
			testRst,
		)
	}
}

func TestIntegration(t *testing.T) {
	const (
		testAddr        = "172.17.0.1:5120" // docker default ip
		testUsername    = "Rick"
		testActName     = "echo"
		testActDirPath  = "./testdata/echo"
		testActArg      = "I'm Mr. Meeseeks! Look at me!"
		testExpectedRst = "I'm Mr. Meeseeks! Look at me!"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := agent.Setup(ctx, agent.Config{
		AccessPoint: testAddr,

		Acts: acts.Config{},
		Pool: pool.Config{Images: []string{"mee6aas/runtime-nodejs"}},
	}); err != nil {
		t.Fatalf("Failed to setup agent: %v", err)
	}

	go func() {
		if err := agent.Serve(ctx, testAddr); err != nil {
			t.Logf("Failed to Serve agent")
		}
	}()
	defer func() {
		if err := agent.Destroy(ctx); err != nil {
			t.Logf("Failed to destory agent: %v", err)
		}
	}()

	invker := mockInvoker.Invoker{}
	defer invker.Close()

	if err := invker.Connect(testAddr); err != nil {
		t.Fatalf("Failed to connect to invoker service: %v", err)
	}

	if err := invker.Register(ctx, testUsername, testActName, testActDirPath); err != nil {
		t.Fatalf("Failed to request for register an activity to agent: %v", err)
	}

	time.Sleep(time.Millisecond * 100)

	if rst, err := invker.Invoke(ctx, testUsername, testActName, testActArg); err == nil {
		if rst != testExpectedRst {
			t.Fatalf("Expected that the result is %s but is %s",
				testExpectedRst,
				rst,
			)
		}
	} else {
		t.Fatalf("Failed to invoke an activity: %v", err)
	}
}
