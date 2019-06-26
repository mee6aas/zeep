package container

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

// Config holds the configuration for the container.
type Config struct {
	Image   string // Name of the image worker use
	Storage string // Path of the directory worker use
}

func init() {
	const (
		apiVersion = dockerEngineAPIVersion
	)

	var (
		err    error
		client *docker.Client
	)

	if client, err = docker.NewClientWithOpts(docker.WithVersion(apiVersion)); err != nil {
		panic(err)
	}

	engineClient = client
}
