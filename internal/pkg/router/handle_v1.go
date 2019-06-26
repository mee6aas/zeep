package router

import (
	"net"

	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"
)

// HandleV1 _
type HandleV1 interface {
	Connected(*net.TCPAddr) error
	Disconnected(*net.TCPAddr)

	SendRequested(*v1.SendRequest, *net.TCPAddr) (*v1.Message, *net.TCPAddr, error)
}
