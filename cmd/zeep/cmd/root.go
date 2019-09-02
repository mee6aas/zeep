package cmd

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/api"
)

var (
	optDebug     bool   // show debugs
	optAgentHost string // host of the agent
	optAgentPort string // port of the agent
	optAgentName string // name of the container that the agent runs
	optUsername  string // username to use for the request.
	optActName   string // name of the activity
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeep",
	Short: "zeep is client for local agent of Mee6aaS",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&optDebug, "debug", false, "print debug logs")
	rootCmd.PersistentFlags().StringVarP(&optAgentName, "agent-name", "N", api.AgentDefaultContainerName, "name of the agent container")
	rootCmd.PersistentFlags().StringVarP(&optAgentHost, "agent-host", "H", "0.0.0.0", "host of the agent serves")
	rootCmd.PersistentFlags().StringVarP(&optAgentPort, "agent-port", "P", strconv.Itoa(api.AgentDefaultPort), "port of the agent serves")

	rootCmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		if optDebug {
			log.SetLevel(log.DebugLevel)
		}
	}
}
