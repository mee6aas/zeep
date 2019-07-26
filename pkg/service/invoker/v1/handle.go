package v1

import "context"

// InvokerAPIServerHandle handles server events.
type InvokerAPIServerHandle interface {
	Requested(context.Context, *InvokeRequest) (*InvokeResponse, error)
}
