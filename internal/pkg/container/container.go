package container

import (
	"context"
	"os"

	dockerCont "github.com/docker/docker/api/types/container"
	dockerMnt "github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"

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
	agentDlgt string
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
func NewContainer(ctx context.Context, conf Config) (c Container, e error) {
	cConf := &dockerCont.Config{
		Image: conf.Image,
		Env: []string{
			api.AgentHostEnvKey + "=" + agentHost,
			api.AgentPortEnvKey + "=" + agentPort,
		},
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd: []string{
			"[ -f " + api.RuntimeSetup + "] && " + api.RuntimeSetup + "; " + api.RuntimeSpawn,
		},
	}

	mnts := []dockerMnt.Mount{
		dockerMnt.Mount{
			Type:   dockerMnt.TypeBind,
			Source: conf.Storage,
			Target: api.ActivityStorage,
			BindOptions: &dockerMnt.BindOptions{
				Propagation: dockerMnt.PropagationShared,
			},
		},
	}

	if agentDlgt != "" {
		cConf.Entrypoint = []string{api.AgentDelegate}
		cConf.Cmd = []string{}

		if log.GetLevel() == log.DebugLevel {
			cConf.Cmd = append(cConf.Cmd, "--debug")
		}

		mnts = append(mnts, dockerMnt.Mount{
			Type:   dockerMnt.TypeBind,
			Source: agentDlgt,
			Target: api.AgentDelegate,
		})
	}

	res, e := engineClient.ContainerCreate(ctx, cConf, &dockerCont.HostConfig{
		NetworkMode: dockerCont.NetworkMode(agentNet),
		Mounts:      mnts,
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
	agentDlgt = os.Getenv(api.AgentDelegateHostPathEnvKey)

	if agentNet == "" {
		agentNet = "bridge"
	}
}
