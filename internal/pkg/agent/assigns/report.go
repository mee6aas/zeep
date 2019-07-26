package assigns

// Report passes result of the invocation to the invoker.
func Report(id string, rst interface{}) (ok bool) {
	var (
		c chan interface{}
	)

	if c, ok = assigns[id]; !ok {
		return
	}

	c <- rst

	close(c)
	delete(assigns, id)

	return
}
