package assigns

// GetAssigneeFromIP gets the username of the assignee that matches the specified address.
func GetAssigneeFromIP(addr string) (u string, ok bool) {
	for _, v := range assignments {
		if v.address == addr {
			u = v.assignee
			ok = true
			return
		}
	}

	return
}
