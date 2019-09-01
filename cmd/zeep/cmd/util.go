package cmd

import (
	"context"

	"github.com/pkg/errors"

	dockerTypes "github.com/docker/docker/api/types"
	dockerFilters "github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"

	"github.com/mee6aas/zeep/api"
)

func createDockerClient() (*docker.Client, error) {
	return docker.NewClientWithOpts(docker.WithVersion(api.DockerAPIVersion))
}

func getAgentAddr() string { return optAgentHost + ":" + optAgentPort }

func getAgentNetworks(ctx context.Context, client *docker.Client) (
	[]dockerTypes.NetworkResource, error,
) {
	return client.NetworkList(ctx, dockerTypes.NetworkListOptions{
		Filters: dockerFilters.NewArgs(
			dockerFilters.KeyValuePair{
				Key:   "name",
				Value: optAgentNet,
			},
		),
	})
}

func getAgentNetwork(ctx context.Context, client *docker.Client) (
	n dockerTypes.NetworkResource, e error,
) {
	ns, e := getAgentNetworks(ctx, client)
	if e != nil {
		return
	}

	if len(ns) > 1 {
		e = errors.Errorf("There are %d networks with the name %s", len(ns), optAgentNet)
		return
	}

	if len(ns) == 0 {
		return
	}

	n = ns[0]

	return
}

func getAgentContainers(ctx context.Context, client *docker.Client) (
	[]dockerTypes.Container, error,
) {
	return client.ContainerList(ctx, dockerTypes.ContainerListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(
			dockerFilters.KeyValuePair{
				Key:   "name",
				Value: optAgentName,
			},
		),
	})
}

func getAgentContainer(ctx context.Context, client *docker.Client) (
	c dockerTypes.Container, e error,
) {
	cs, e := getAgentContainers(ctx, client)
	if e != nil {
		return
	}

	if len(cs) > 1 {
		e = errors.Errorf("There are %d containers with the name %s", len(cs), optAgentName)
		return
	}

	if len(cs) == 0 {
		return
	}

	c = cs[0]

	return
}
