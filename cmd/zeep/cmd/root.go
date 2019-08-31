package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// options
var (
	agentAddr string // address to request. (default: )
	username  string // username to use for request. (default: jerry)
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
}
