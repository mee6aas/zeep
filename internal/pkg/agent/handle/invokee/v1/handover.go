package v1

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/internal/pkg/agent/assistant/handover"
)

// HandoverRequested is invoked when the invokee client requested to handover its control to other.
func (h Handle) HandoverRequested(addr *net.TCPAddr) (e error) {
	l := log.WithFields(log.Fields{
		"addr": addr.String(),
	})

	ip := addr.IP.String()

	if ok := handover.Reserve(ip); !ok {
		l.Warn("Duplicated request for handover")
		e = errors.New("Duplicated request for handover")
		return
	}

	return
}
