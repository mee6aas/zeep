package cmd

import (
	"github.com/mee6aas/zeep/api"
	"github.com/spf13/cobra"
)

var (
	optAgentNet string // name of the network that agent serves
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage the agent",
}

func init() {
	rootCmd.AddCommand(agentCmd)

	agentCmd.PersistentFlags().StringVar(&optAgentNet, "agent-net", api.AgentDefaultNetworkName, "name of the network that the agent serves")
}
