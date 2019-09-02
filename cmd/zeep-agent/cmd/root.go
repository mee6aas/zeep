package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/internal/pkg/agent"
	"github.com/mee6aas/zeep/internal/pkg/worker/pool"
)

var (
	optDebug  bool
	agentAddr string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeep-agent",
	Short: "zeep-agent is agent of Mee6aaS",

	RunE: func(cmd *cobra.Command, args []string) (e error) {
		{
			log.Info("Setting up agent")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			e = agent.Setup(ctx, agent.Config{
				Pool: pool.Config{Images: []string{"mee6aas/runtime-nodejs:latest"}},
			})
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to setup agent")
				return
			}

			log.Info("Agent setup")
		}

		log.WithFields(log.Fields{
			"addr": agentAddr,
		}).Info("Serving agent")

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGINT)
			<-c

			{
				log.Info("Destroying agent")

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
				e = agent.Destroy(ctx)
				cancel()

				if e != nil {
					log.WithError(e).Error("Failed to destroy agent")
					return
				}
			}

			wg.Done()
		}()

		if e = agent.Serve(context.Background(), agentAddr); e != nil {
			log.WithError(e).Error("Failed to serve agent")
		}

		log.Info("Agent destroyed")

		wg.Wait()

		return
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
	rootCmd.PersistentFlags().BoolVar(&optDebug, "debug", false, "print debug messages")
	rootCmd.Flags().StringVar(&agentAddr, "addr", "0.0.0.0:5122", "address to serve")

	rootCmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		if optDebug {
			log.SetLevel(log.DebugLevel)
		}
	}
}
