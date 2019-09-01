package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
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
			log.Info("Setting up agent")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			err := agent.Setup(ctx, agent.Config{
				Addr: agentAddr,
				Pool: pool.Config{Images: []string{"mee6aas/runtime-nodejs:latest"}},
			})
			cancel()

			if err != nil {
				log.WithError(err).Error("Failed to setup agent")
				return
			}

			log.Info("Agent setup")
		}

		{
			log.WithFields(log.Fields{
				"addr": agentAddr,
			}).Info("Serving agent")

			if err := agent.Serve(context.Background()); err != nil {
				log.WithError(err).Error("Failed to serve agent")
			}

			log.Info("Agent stopped")
		}

		{
			log.Info("Destroying agent")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			err := agent.Destroy(ctx)
			cancel()

			if err != nil {
				log.WithError(err).Error("Failed to destroy agent")
			}

			log.Info("Agent destroyed")
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
