package v1

// InvokerAPIServerHandle handles server events.
type InvokerAPIServerHandle interface {
	Requested(*InvokeRequest) (*InvokeResponse, error)
}
