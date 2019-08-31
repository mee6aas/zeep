package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	docker "github.com/docker/docker/client"
	"github.com/mee6aas/zeep/api"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			client *docker.Client
			contID string
			netID  string
			tmpDir string
		)

		{
			fmt.Print("Create Docker client...")
			client, err = docker.NewClientWithOpts(docker.WithVersion(api.DockerAPIVersion))
			if err != nil {
				fmt.Println("Faile")
				fmt.Println(err)
				return
			}
			fmt.Println("Done")
		}

		{
			fmt.Print("Check if the container for agent is exists...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, err := getAgentContainers(ctx, client)
			cancel()

			if err != nil {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}

			if len(res) > 1 {
				fmt.Println("Failed")
				fmt.Printf("There are %d containers with the name %s", len(res), api.AgentDefaultContainerName)
				return
			}

			if len(res) == 1 {
				contID = res[0].ID
				fmt.Println("Found")
			}

			if len(res) == 0 {
				fmt.Println("Not found")
			}
		}

		if contID != "" {
			{
				fmt.Printf("Stop container %s...", contID)
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				err = client.ContainerStop(ctx, contID, nil)
				cancel()

				if err != nil {
					fmt.Println("Failed")
					fmt.Println(err)
					return
				}

				fmt.Println("Done")
			}

			{
				fmt.Printf("Inspect container %s...", contID)
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				res, err := client.ContainerInspect(ctx, contID)
				cancel()

				if err != nil {
					fmt.Println("Failed")
					fmt.Println(err)
					return
				}
				fmt.Println("Done")

				fmt.Printf("Delete temp directory for agent...")

				for _, mnt := range res.Mounts {
					if mnt.Destination != "/tmp" {
						continue
					}

					tmpDir = mnt.Source

					break
				}

				if tmpDir == "" {
					fmt.Println("Not found")
				} else {
					if err = os.RemoveAll(tmpDir); err == nil {
						fmt.Println("Done")
					} else {
						fmt.Println("Failed")
						fmt.Println(err)
					}
				}
			}
		}

		{
			fmt.Print("Check if the network for agent is exists...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, err := getAgentNetworks(ctx, client)
			cancel()

			if err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}

			if len(res) > 1 {
				fmt.Println("Failed")
				fmt.Printf("There are %d networks with the name %s", len(res), api.AgentDefaultNetworkName)
				return
			}

			if len(res) == 1 {
				netID = res[0].ID
				fmt.Println("Found")
			}

			if len(res) == 0 {
				fmt.Println("Not found")
			}
		}

		if netID != "" {
			fmt.Printf("Remove network %s...", netID)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			err = client.NetworkRemove(ctx, netID)
			cancel()

			if err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}

			fmt.Println("Done")
		}
	},
}

func init() {
	agentCmd.AddCommand(destroyCmd)
}
