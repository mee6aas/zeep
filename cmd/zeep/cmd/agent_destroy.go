package cmd

import (
	"context"
	"os"
	"time"

	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",

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
			log.Info("Docker client created")
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
				log.WithField("ID", contID).Info("Container that agent runs decided")
			}
		}

		if contID != "" {
			{
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				e = client.ContainerStop(ctx, contID, nil)
				cancel()

				if e != nil {
					log.WithError(e).Error("Failed to stop container")
					return
				}

				log.Info("Container stopped")
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
						log.Info("Temp direcotry for agent removed")
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
				log.WithField("ID", netID).Info("Network that agent served decided")
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

			log.Info("Network removed")
		}

		return
	},
}

func init() {
	agentCmd.AddCommand(destroyCmd)
}
