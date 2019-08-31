package cmd

import (
	"context"

	dockerTypes "github.com/docker/docker/api/types"
	dockerFilters "github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/api"
)

var (
	// name of the network that agent uses
	agentNetName string

	// name of the container that the agent runs
	agentContName string
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage the agent",
}

func getAgentNetworks(ctx context.Context, client *docker.Client) (
	[]dockerTypes.NetworkResource, error,
) {
	return client.NetworkList(ctx, dockerTypes.NetworkListOptions{
		Filters: dockerFilters.NewArgs(
			dockerFilters.KeyValuePair{
				Key:   "name",
				Value: agentNetName,
			},
		),
	})
}

func getAgentContainers(ctx context.Context, client *docker.Client) (
	[]dockerTypes.Container, error,
) {
	return client.ContainerList(ctx, dockerTypes.ContainerListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(
			dockerFilters.KeyValuePair{
				Key:   "name",
				Value: agentContName,
			},
		),
	})
}

func init() {
	rootCmd.AddCommand(agentCmd)

	agentCmd.PersistentFlags().StringVar(&agentNetName, "net", api.AgentDefaultNetworkName, "name of the network that agent serves")
	agentCmd.PersistentFlags().StringVar(&agentContName, "name", api.AgentDefaultContainerName, "name of the container that agent runs")
}
