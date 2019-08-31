package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	dockerCont "github.com/docker/docker/api/types/container"
	dockerMnt "github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/api"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve agent in docker container",

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
			if client, err = docker.NewClientWithOpts(docker.WithVersion(api.DockerAPIVersion)); err != nil {
				fmt.Println("Faile")
				fmt.Println(err)
				return
			}
			fmt.Println("Done")
		}

		{
			fmt.Print("Check if the network for agent is exists...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			if res, err := getAgentNetworks(ctx, client); err == nil {
				cancel()
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
			} else {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}
		}

		if netID == "" {
			fmt.Print("Create network for agent...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			if res, err := client.NetworkCreate(ctx, agentNetName, dockerTypes.NetworkCreate{
				Driver: "bridge",
				Scope:  "local",
			}); err == nil {
				netID = res.ID
				cancel()
				fmt.Println("Done")
				fmt.Println(netID)
			} else {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}
		}

		{
			fmt.Print("Check if the container for agent is exists...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			if res, err := getAgentContainers(ctx, client); err == nil {
				cancel()
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
			} else {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}
		}

		if contID == "" {
			fmt.Print("Create container for agent...")

			if tmpDir, err = ioutil.TempDir("", ""); err != nil {
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			if res, err := client.ContainerCreate(ctx, &dockerCont.Config{
				Image: "mee6aas/zeep:latest",
				Env: []string{
					api.AgentTmpDirPathEnvKey + "=" + tmpDir,
					api.AgentNetworkEnvKey + "=" + agentNetName,
					api.AgentHostEnvKey + "=" + agentContName,
					api.AgentPortEnvKey + "=" + strconv.Itoa(api.AgentDefaultPort),
				},
			}, &dockerCont.HostConfig{
				Privileged:  true,
				NetworkMode: dockerCont.NetworkMode(agentNetName),
				Mounts: []dockerMnt.Mount{
					dockerMnt.Mount{
						Type:   dockerMnt.TypeBind,
						Source: tmpDir,
						Target: "/tmp",
					},
					dockerMnt.Mount{
						Type:   dockerMnt.TypeBind,
						Source: "/var/run/docker.sock",
						Target: "/var/run/docker.sock",
					},
				},
			}, nil, api.AgentDefaultContainerName); err == nil {
				contID = res.ID
				cancel()
				fmt.Println("Done")
				fmt.Println(contID)
			} else {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}
		}

		{
			fmt.Print("Start agent...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			if err = client.ContainerStart(ctx, contID, dockerTypes.ContainerStartOptions{}); err == nil {
				cancel()
				fmt.Println("Done")
			} else {
				cancel()
				fmt.Println("Failed")
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	agentCmd.AddCommand(serveCmd)
}
