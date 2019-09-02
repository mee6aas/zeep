package cmd

import (
	"github.com/spf13/cobra"
)

// actCmd represents the act command
var actCmd = &cobra.Command{
	Use:   "act",
	Short: "Manage activities",
}

func init() {
	rootCmd.AddCommand(actCmd)
}
