package assigns

// GetAssigneeFromIP returns username who invokes an activity processed in given address.
func GetAssigneeFromIP(addr string) (u string, ok bool) {
	for _, v := range assigns {
		if v.address == addr {
			u = v.assignee
			ok = true
			return
		}
	}

	return
}
