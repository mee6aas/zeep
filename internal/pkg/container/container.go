package container

import (
	"context"
	"os"

	dockerCont "github.com/docker/docker/api/types/container"
	dockerMnt "github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"

	"github.com/mee6aas/zeep/api"
)

const (
	dockerEngineAPIVersion = api.DockerAPIVersion
)

var (
	engineClient *docker.Client

	agentNet  string
	agentHost string
	agentPort string
)

// Container describes container.
type Container struct {
	id string // ID of container
	ip string // IP of container
}

// ID returns the ID of this container.
func (c Container) ID() string { return c.id }

// IP returns the IP address of this container.
func (c Container) IP() string { return c.ip }

// Config holds the configuration for the container.
type Config struct {
	Image   string // Image to create container
	Storage string // Path of the directory to be mounted on container
}

// NewContainer creates a new container based on the given configuration
// and returns its descriptor.
func NewContainer(ctx context.Context, config Config) (c Container, e error) {
	res, e := engineClient.ContainerCreate(ctx, &dockerCont.Config{
		Image: config.Image,
		Env: []string{
			api.AgentHostEnvKey + "=" + agentHost,
			api.AgentPortEnvKey + "=" + agentPort,
		},
	}, &dockerCont.HostConfig{
		NetworkMode: dockerCont.NetworkMode(agentNet),
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
	}, nil, "")
	if e != nil {
		return c, e
	}

	c.id = res.ID

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

	agentNet = os.Getenv(api.AgentNetworkEnvKey)
	agentHost = os.Getenv(api.AgentHostEnvKey)
	agentPort = os.Getenv(api.AgentPortEnvKey)

	if agentNet == "" {
		agentNet = "bridge"
	}
}
