package cmd

import (
	"context"
	"time"

	"github.com/pkg/errors"

	dockerCont "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// agentStopCmd represents the stop command
var agentStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop agent",
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			client *docker.Client
			contID string
		)

		{
			if client, e = createDockerClient(); e != nil {
				log.WithError(e).Error("Failed to create Docker client")
				return
			}
			log.Debug("Docker client created")
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			c, e := getAgentContainer(ctx, client)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to decide client of the agent")
				return e
			}

			if c.ID != "" {
				contID = c.ID
				log.WithField("ID", contID).Debug("Container that agent runs decided")
			}
		}

		if contID == "" {
			e = errors.New("Not found")
			log.Error("There is no agent container")
			return
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			// TODO: is it need?
			timeout := time.Minute
			e = client.ContainerStop(ctx, contID, &timeout)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to stop container")
				return
			}

			log.Debug("Container stop requested")
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			cB, cE := client.ContainerWait(ctx, contID, dockerCont.WaitConditionNotRunning)

			select {
			case b := <-cB:
				log.Debug(b)
			case e = <-cE:
			}

			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to wait stop agent")
				return
			}

			log.Debug("Container stopped")
		}

		return
	},
}

func init() {
	agentCmd.AddCommand(agentStopCmd)

	agentStopCmd.Flags().BoolVar(&optDetach, "detach", false, "don't wait container stopped")
}
