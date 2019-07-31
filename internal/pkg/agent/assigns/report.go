package assigns

// Report passes result of the invocation to the invoker.
func Report(invkID string, rst interface{}) (ok bool) {
	var (
		c chan interface{}
	)

	if c, ok = assigns[invkID]; !ok {
		return
	}

	c <- rst

	close(c)
	delete(assigns, invkID)

	return
}
