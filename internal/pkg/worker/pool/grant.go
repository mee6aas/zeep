package pool

import "github.com/mee6aas/zeep/internal/pkg/worker"

// Grant move worker in pending list to ready list.
func (p *Pool) Grant(ip string, ta worker.TaskAssigner, version string) (ok bool) {
	var (
		w worker.Worker
	)

	if w, ok = p.pendings[ip]; !ok {
		return
	}

	w.InvokeeVersion = version
	w.Allocate(ta)

	go func() {
		img := w.Image()
		p.granted[img] <- w
		delete(p.pendings, ip)
	}()

	return
}
