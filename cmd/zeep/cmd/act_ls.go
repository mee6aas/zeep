package cmd

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	invokerAPI "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// actLsCmd represents the ls command
var actLsCmd = &cobra.Command{
	Use:   "ls USERNAME",
	Short: "List activities",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			username = optUsername
		)

		if len(args) == 1 {
			username = args[0]
		}

		// make request
		{
			l := log.WithFields(log.Fields{
				"addr": getAgentAddr(),
				"user": username,
			})

			conn, e := grpc.Dial(getAgentAddr(), grpc.WithInsecure())
			if e != nil {
				l.Error("Failed to dial")
				return e
			}
			client := invokerAPI.NewInvokerClient(conn)

			req := &invokerAPI.ListRequest{
				Username: username,
			}

			log.Debug(req)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, e := client.List(ctx, req)
			cancel()

			// TODO: show more information
			if e != nil {
				l.Error("Failed to list activities")
				return e
			}

			for _, a := range res.Activities {
				fmt.Println(a)
			}
		}

		return
	},
}

func init() {
	actCmd.AddCommand(actLsCmd)
}
