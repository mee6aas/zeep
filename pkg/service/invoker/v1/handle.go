package v1

import "context"

// InvokerAPIServerHandle handles server events.
type InvokerAPIServerHandle interface {
	InvokeRequested(context.Context, string, string, string) (*InvokeResponse, error)
	RegisterRequested(context.Context, string, string) error
}
