package container_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	. "github.com/mee6aas/zeep/internal/pkg/container"

	docker "github.com/docker/docker/client"
)

var (
	testContainerCreateFailed = false

	testEngineClient         *docker.Client
	testContainer            Container
	testContainerImage       = "mee6aas/runtime-test:latest"
	testContainerStoragePath string
)

func TestContainerIsExistsFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	c := Container{}

	if c.IsExists(ctx) {
		t.Fatal("Expected to not exists")
	}
}

func TestContainerCreateFail(t *testing.T) {
	var (
		err error
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if testContainer, err = NewContainer(ctx, Config{
		Image:   "not exists",
		Storage: testContainerStoragePath,
	}); err == nil {
		t.Fatal("Expected to fail to create container")
	}
}

func TestContainerCreate(t *testing.T) {
	var (
		err error
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if testContainer, err = NewContainer(ctx, Config{
		Image:   testContainerImage,
		Storage: testContainerStoragePath,
	}); err != nil {
		testContainerCreateFailed = true
		t.Fatalf("Failed to create container %v", err)
	}

	if ok := testContainer.IsExists(ctx); !ok {
		t.Fatalf("Expected to container exists.")
	}
}

func TestContainerRemove(t *testing.T) {
	var (
		err    error
		ctx    = context.Background()
		cancel context.CancelFunc
	)

	if testContainerCreateFailed {
		t.Skipf("TestContainerCreate failed")
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err = testContainer.Remove(ctx); err != nil {
		t.Fatalf("Failed to remove container %v", err)
	}

	if err = os.Remove(testContainerStoragePath); err != nil {
		t.Logf("Failed to remove testing storage %s", testContainerStoragePath)
	}
}

func init() {
	const (
		apiVersion = "1.39"
	)

	var (
		err    error
		client *docker.Client
	)

	if client, err = docker.NewClientWithOpts(docker.WithVersion(apiVersion)); err != nil {
		panic(err)
	}

	testEngineClient = client

	if testContainerStoragePath, err = ioutil.TempDir("", ""); err != nil {
		panic(err)
	}
}
