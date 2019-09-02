package v1

import (
	"context"

	"github.com/mee6aas/zeep/pkg/activity"
)

// InvokerAPIServerHandle handles server events.
type InvokerAPIServerHandle interface {

	// ctx, username, actNAme, arg
	InvokeRequested(context.Context, string, string, string) (*InvokeResponse, error)

	// ctx, username, actName, actDirPath
	AddRequested(context.Context, string, string, string) error

	// ctx, username
	ListRequested(context.Context, string) ([]activity.Activity, error)

	// ctx, username, actName
	RemoveRequested(context.Context, string, string) error
}
