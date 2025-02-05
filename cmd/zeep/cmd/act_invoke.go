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

// actInvokeCmd represents the invoke command
var actInvokeCmd = &cobra.Command{
	Use:   "invoke ACTIVITY_NAME [ARGUMENT]",
	Short: "Invoke activity",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) (e error) {
		var (
			trg = args[0] // name of the activity to invoke
			arg string    // argument to be passed to activity
		)

		if len(args) == 2 {
			arg = args[1]
		}

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

			req := &invokerAPI.InvokeRequest{
				Username: optUsername,
				ActName:  trg,
				Arg:      arg,
			}

			log.Debug(req)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			res, e := client.Invoke(ctx, req)
			cancel()

			// TODO: show more information
			if e != nil {
				l.Error("Failed to invoke activity")
				return e
			}

			fmt.Println(res.GetResult())
		}

		return
	},
}

func init() {
	actCmd.AddCommand(actInvokeCmd)

	actInvokeCmd.Flags().StringVarP(&optUsername, "username", "u", "Jerry", "username to use for request")
}
