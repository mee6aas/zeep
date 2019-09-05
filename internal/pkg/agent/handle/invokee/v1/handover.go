package v1

import (
	"net"
)

// HandoverRequested is invoked when the invokee client requested to handover its control to other.
func (h Handle) HandoverRequested(addr *net.TCPAddr) (e error) {

	return
}
