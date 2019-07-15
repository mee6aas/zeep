package v1

import (
	"context"
	"net"
)

// InvokeeAPIServerHandle hadles server events.
type InvokeeAPIServerHandle interface {
	Connected(context.Context, *net.TCPAddr, chan<- Task) error
	Disconnected(*net.TCPAddr)
	Reported(*ReportRequest)
}
