package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	invokerAPI "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add PATH|URL",
	Short: "Add activity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			src = args[0] // path of the activity to add
			dst string    // path in the host where the activity to be copied temporary.
		)

		// TODO: support URL
		if strings.Contains(src, "://") {
			log.Error("Not implemented for URL input")
			return errors.New("Not implemented")
		}

		// if name not specified, decide it from the source name.
		// e.g. /path/to/activity -> activity
		if optActName == "" {
			optActName = filepath.Base(src)
		}

		// find mountpoint that mounted on `/tmp` in the agent continaer
		{
			var mp string // mountpoint

			client, e := createDockerClient()
			if e != nil {
				log.Error("Failed to create Docker client")
				return e
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			c, e := getAgentContainer(ctx, client)
			cancel()

			if e != nil {
				log.Error("Failed to decide the agent container")
				return e
			}

			if c.ID == "" {
				log.Error("There is no agent container")
				return errors.New("Not found")
			}

			for _, mnt := range c.Mounts {
				if mnt.Destination != "/tmp" {
					continue
				}

				mp = mnt.Source
			}

			if mp == "" {
				log.WithField("ID", c.ID).Error(`Failed to find mountpoint that mounted on "/tmp" in the agent container.`)
				return errors.New("Failed to find mountpoint")
			}

			log.WithField("mountpoint", mp).Debug("Found mountpoint")

			// create temporal directory to copy the activity into the agent container
			if dst, e = ioutil.TempDir(mp, "zeep-act-add-"); e != nil {
				log.Error("Failed to create temp directory")
				return e
			}
		}

		// remove temp directory after the command finished
		defer func() {
			if e := os.RemoveAll(dst); e != nil {
				log.WithError(e).Warn("Filaed to remove temp directory")
			}
		}()

		log.WithFields(log.Fields{
			"src": src,
			"dst": dst,
		}).Debug("Copying activity")

		// copy activity from the host into the container
		if e = copy.Copy(src, dst); e != nil {
			log.WithFields(log.Fields{
				"src": src,
				"dst": dst,
			}).Error("Failed to copy")
			return e
		}

		// make request
		{
			// path in the agent container where the activity to be copied.
			trg := filepath.Join("/tmp", filepath.Base(dst))

			l := log.WithFields(log.Fields{
				"addr": getAgentAddr(),
				"user": optUsername,
				"act":  optActName,
			})

			conn, e := grpc.Dial(getAgentAddr(), grpc.WithInsecure())
			if e != nil {
				l.Error("Failed to dial")
				return e
			}
			client := invokerAPI.NewInvokerClient(conn)

			req := &invokerAPI.RegisterRequest{
				Username: optUsername,
				ActName:  optActName,
				Method:   invokerAPI.RegisterMethod_LOCAL,
				Path:     trg,
			}

			log.Debug(req)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			_, e = client.Register(ctx, req)
			cancel()

			// TODO: show more information
			if e != nil {
				l.Error("Failed to add activity")
				return e
			}
		}

		log.WithFields(log.Fields{
			"user": optUsername,
			"name": optActName,
		}).Info("Activity registered")

		return
	},
}

func init() {
	actCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&optActName, "name", "n", "", "name of the activity to add")
}
