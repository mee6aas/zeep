package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/internal/pkg/agent"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	agentAddr string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeep-agent",
	Short: "zeep-agent is agent of Mee6aaS",

	Run: func(cmd *cobra.Command, args []string) {
		{
			fmt.Print("Setup agent...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			if err := agent.Setup(ctx, agent.Config{
				Addr: agentAddr,
				Pool: pool.Config{Images: []string{"mee6aas/runtime-nodejs:latest"}},
			}); err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
				cancel()
				return
			}
			cancel()
			fmt.Println("Done")
		}

		{
			fmt.Printf("Serve agent at [%s]...", agentAddr)
			if err := agent.Serve(context.Background()); err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
			}
			fmt.Println("Stopped")
		}

		{
			fmt.Print("Destroy agent...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			if err := agent.Destroy(ctx); err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
			}
			cancel()
			fmt.Println("Done")
		}
	},
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
	rootCmd.Flags().StringVar(&agentAddr, "addr", "0.0.0.0:5122", "address to serve")
}
