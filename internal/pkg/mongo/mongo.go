package mongo

import (
	docker "github.com/docker/docker/client"
	"github.com/mee6aas/zeep/api"
)

const (
	dockerEngineAPIVersion = api.DockerAPIVersion
)

var (
	engineClient *docker.Client
)

func init() {
	const (
		apiVersion = dockerEngineAPIVersion
	)

	var (
		err    error
		client *docker.Client
	)

	client, err = docker.NewClientWithOpts(docker.WithVersion(apiVersion))
	if err != nil {
		panic(err)
	}

	engineClient = client
}
