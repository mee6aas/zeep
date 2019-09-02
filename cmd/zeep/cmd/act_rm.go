package cmd

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	invokerAPI "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// actRmCmd represents the rm command
var actRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove activity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			trg = args[0] // name of the activity to invoke
		)

		// make request
		{
			l := log.WithFields(log.Fields{
				"addr": getAgentAddr(),
				"user": optUsername,
				"act":  trg,
			})

			conn, e := grpc.Dial(getAgentAddr(), grpc.WithInsecure())
			if e != nil {
				l.Error("Failed to dial")
				return e
			}
			client := invokerAPI.NewInvokerClient(conn)

			req := &invokerAPI.RemoveRequest{
				Username: optUsername,
				ActName:  trg,
			}

			log.Debug(req)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			_, e = client.Remove(ctx, req)
			cancel()

			// TODO: show more information
			if e != nil {
				l.Error("Failed to invoke activity")
				return e
			}
		}

		return
	},
}

func init() {
	actCmd.AddCommand(actRmCmd)

	actRmCmd.Flags().StringVarP(&optUsername, "username", "u", "Jerry", "username to use for request")
}
