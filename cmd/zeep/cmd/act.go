package cmd

import (
	"strconv"

	"github.com/mee6aas/zeep/api"
	"github.com/spf13/cobra"
)

// actCmd represents the act command
var actCmd = &cobra.Command{
	Use:   "act",
	Short: "Manage activities",
}

func init() {
	rootCmd.AddCommand(actCmd)

	actCmd.PersistentFlags().StringVarP(&optUsername, "username", "u", "Jerry", "username to use for request")
	actCmd.PersistentFlags().StringVarP(&optAgentHost, "agent-host", "H", "0.0.0.0", "host of the agent to request")
	actCmd.PersistentFlags().StringVarP(&optAgentPort, "agent-port", "P", strconv.Itoa(api.AgentDefaultPort), "port of the agent to request")
}
