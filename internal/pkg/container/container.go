package container

import (
	"context"

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
	image   string // Image used to create container
	storage string // Path of the directory mounted on container
}

// ID returns the ID of this container.
func (c *Container) ID() string { return c.id }

// Image returns the image used to create this container.
func (c *Container) Image() string { return c.image }

// Storage returns the path of the directory that mounted on this container.
func (c *Container) Storage() string { return c.storage }

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
		Env:   []string{},
	}, &dockerCont.HostConfig{
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
