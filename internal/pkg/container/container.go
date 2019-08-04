package container

import (
	"context"
	"os"

	"github.com/mee6aas/zeep/api"

	dockerCont "github.com/docker/docker/api/types/container"
	dockerMnt "github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
)

const (
	dockerEngineAPIVersion = api.DockerAPIVersion
)

var (
	engineClient *docker.Client
)

// Container describes container.
type Container struct {
	id      string // ID of container
	ip      string // TCP address of container
	image   string // Image used to create container
	storage string // Path of the directory mounted on container
}

// ID returns the ID of this container.
func (c Container) ID() string { return c.id }

// IP returns the IP address of this container.
func (c Container) IP() string { return c.ip }

// Image returns the image used to create this container.
func (c Container) Image() string { return c.image }

// Storage returns the path of the directory that mounted on this container.
func (c Container) Storage() string { return c.storage }

// Config holds the configuration for the container.
type Config struct {
	Image   string // Image to create container
	Storage string // Path of the directory to be mounted on container
}

// NewContainer creates a new container based on the given configuration
// and returns its descriptor.
func NewContainer(ctx context.Context, config Config) (cont Container, err error) {
	var (
		res dockerCont.ContainerCreateCreatedBody
	)

	if res, err = engineClient.ContainerCreate(ctx, &dockerCont.Config{
		Image: config.Image,
		Env: []string{
			"AGENT_HOST=" + os.Getenv("AGENT_HOST"),
			"AGENT_PORT=" + os.Getenv("AGENT_PORT"),
		},
	}, &dockerCont.HostConfig{
		NetworkMode: api.NetworkName,
		Mounts: []dockerMnt.Mount{
			dockerMnt.Mount{
				Type:   dockerMnt.TypeBind,
				Source: config.Storage,
				Target: api.ActivityStorage,
				BindOptions: &dockerMnt.BindOptions{
					Propagation: dockerMnt.PropagationShared,
				},
			},
		},
	}, nil, ""); err != nil {
		return
	}

	cont.id = res.ID
	cont.image = config.Image
	cont.storage = config.Storage

	return
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
