package cmd

import (
	"context"
	"os"
	"time"

	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// agentDestroyCmd represents the destroy command
var agentDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destory agent",
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			client *docker.Client
			contID string
			netID  string
			tmpDir string
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

		if contID != "" {
			{
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				timeout := time.Minute
				// TODO: remove?
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
				res, e := client.ContainerInspect(ctx, contID)
				cancel()

				if e != nil {
					log.WithError(e).Error("Failed to inspect container")
					return e
				}

				for _, mnt := range res.Mounts {
					if mnt.Destination != "/tmp" {
						continue
					}

					tmpDir = mnt.Source

					break
				}

				if tmpDir == "" {
					log.Warn("Temp direcotry for agent not found")
				} else {
					if e = os.RemoveAll(tmpDir); e == nil {
						log.Debug("Temp direcotry for agent removed")
					} else {
						log.WithError(e).Warn("Failed to remove temp direcotry for agent")
					}
				}
			}
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			n, e := getAgentNetwork(ctx, client)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to decide network of the agent")
				return e
			}

			if n.ID != "" {
				netID = n.ID
				log.WithField("ID", netID).Debug("Network that agent served decided")
			}
		}

		if netID != "" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			e = client.NetworkRemove(ctx, netID)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to remove network")
				return
			}

			log.Debug("Network removed")
		}

		return
	},
}

func init() {
	agentCmd.AddCommand(agentDestroyCmd)
}
