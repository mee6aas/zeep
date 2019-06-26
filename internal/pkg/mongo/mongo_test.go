package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/mee6aas/zeep/api"
	. "github.com/mee6aas/zeep/internal/pkg/mongo"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

var (
	testEngineClient *docker.Client
	testMongoVersion = api.DefaultMongoVersion
	testImageExists  = false
)

//
// Remove image mongo:{testMongoVersion} to full test.
//

func TestDeploy(t *testing.T) {
	var (
		err    error
		res    DeployDeployedBody
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if res, err = Deploy(ctx, WithVersion("notExists")); err == nil {
		t.Fatalf("Expected to fail to deploy")
	}

	if !testImageExists {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if res, err = Deploy(ctx, WithVersion(testMongoVersion)); err == nil {
			t.Fatalf("Expected to fail to deploy")
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err = Deploy(ctx, WithVersion(testMongoVersion), WithPull()); err == nil {
		t.Fatalf("Expected to fail to deploy")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err = Deploy(ctx, WithVersion(testMongoVersion), WithNetwork("host")); err != nil {
		t.Fatalf("Failed to deploy: %v", err)
	}

	_ = res
}

func init() {
	const (
		apiVersion = "1.39"
	)

	var (
		err    error
		ctx    context.Context
		cancel context.CancelFunc
		client *docker.Client
	)

	if client, err = docker.NewClientWithOpts(docker.WithVersion(apiVersion)); err != nil {
		panic(err)
	}

	testEngineClient = client

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	f := filters.NewArgs()
	f.Add("reference", "mongo:"+testMongoVersion)

	if imgSum, e := client.ImageList(ctx, types.ImageListOptions{
		Filters: f,
	}); e == nil {
		testImageExists = len(imgSum) != 0
	} else {
		panic(e)
	}
}