package assigns

// Report passes result of the invocation to the invoker.
func Report(invkID string, rst interface{}) (ok bool) {
	var (
		a assign
	)

	if a, ok = assigns[invkID]; !ok {
		return
	}

	a.holder <- rst

	close(a.holder)
	delete(assigns, invkID)

	return
}
