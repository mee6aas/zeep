package v1

import "context"

// InvokerAPIServerHandle handles server events.
type InvokerAPIServerHandle interface {
	InvokeRequested(context.Context, string, string, string, string) (*InvokeResponse, error)

	// ctx, username, actName, actDirPath
	RegisterRequested(context.Context, string, string, string) error
}
