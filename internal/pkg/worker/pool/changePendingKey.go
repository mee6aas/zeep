// +build int

package pool

// ChangePendingKey changes key mapped to worker.
// This method must be used only for the test.
func (p *Pool) ChangePendingKey(to string, from string) (ok bool) {
	if _, ok = p.pendings[from]; !ok {
		return
	}

	p.pendings[to] = p.pendings[from]
	delete(p.pendings, from)

	return
}
