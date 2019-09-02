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
	dockerNat "github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mee6aas/zeep/api"
)

// agentServeCmd represents the serve command
var agentServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve agent in docker container",

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
			n, e := getAgentNetwork(ctx, client)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to decide network of the agent")
				return e
			}

			if n.ID != "" {
				netID = n.ID
				log.Debug("Use existing network")
			}
		}

		if netID == "" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, e := client.NetworkCreate(ctx, optAgentNet, dockerTypes.NetworkCreate{
				Driver: "bridge",
				Scope:  "local",
			})
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to create network")
				return e
			}

			netID = res.ID
			log.Debug("Network created")
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
				log.Debug("Use existing container")
			}
		}

		if contID == "" {
			if tmpDir, e = ioutil.TempDir("", "zeep-cli-"); e != nil {
				log.WithError(e).Error("Failed to create temp direcotry")
				return
			}

			var cmd []string
			if optDebug {
				cmd = []string{"--debug"}
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, e := client.ContainerCreate(ctx, &dockerCont.Config{
				StopSignal: "SIGINT",
				Cmd:        cmd,
				Image:      "mee6aas/zeep:latest",
				Env: []string{
					api.AgentTmpDirPathEnvKey + "=" + tmpDir,
					api.AgentNetworkEnvKey + "=" + optAgentNet,
					api.AgentHostEnvKey + "=" + optAgentName,
					api.AgentPortEnvKey + "=" + strconv.Itoa(api.AgentDefaultPort),
				},
			}, &dockerCont.HostConfig{
				Privileged:  true,
				NetworkMode: dockerCont.NetworkMode(optAgentNet),
				PortBindings: dockerNat.PortMap{
					dockerNat.Port(optAgentPort + "/tcp"): []dockerNat.PortBinding{
						dockerNat.PortBinding{
							HostIP:   optAgentHost,
							HostPort: optAgentPort,
						},
					},
				},
				Mounts: []dockerMnt.Mount{
					dockerMnt.Mount{
						Type:   dockerMnt.TypeBind,
						Source: tmpDir,
						Target: "/tmp",
						BindOptions: &dockerMnt.BindOptions{
							Propagation: dockerMnt.PropagationShared,
						},
					},
					dockerMnt.Mount{
						Type:   dockerMnt.TypeBind,
						Source: "/var/run/docker.sock",
						Target: "/var/run/docker.sock",
					},
				},
			}, nil, api.AgentDefaultContainerName)
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to create container")
				return e
			}

			contID = res.ID

			log.WithField("ID", contID).Debug("Container created")
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			e = client.ContainerStart(ctx, contID, dockerTypes.ContainerStartOptions{})
			cancel()

			if e != nil {
				log.WithError(e).Error("Failed to start container")
				return
			}

			log.Debug("Agent started")
		}

		fmt.Println(contID)

		return
	},
}

func init() {
	agentCmd.AddCommand(agentServeCmd)
}
