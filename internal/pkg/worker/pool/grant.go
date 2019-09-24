package pool

import (
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/internal/pkg/worker"
)

// Grant move worker in pending list to ready list.
func (p *Pool) Grant(ip string, ta worker.TaskAssigner, version string) bool {
	var (
		w worker.Worker
	)

	w, ok := p.pendings[ip]
	if !ok {
		log.WithField("IP", ip).Warn("Not exists in the pended list")
		return false
	}

	w.InvokeeVersion = version
	w.Allocate(ta)

	go func() {
		img := w.Image()
		p.granted[img] <- w
		delete(p.pendings, ip)
	}()

	return true
}
