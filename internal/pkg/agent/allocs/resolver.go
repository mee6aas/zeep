package allocs

// ResolveNameFromIP returns name mapped to given IP.
func ResolveNameFromIP(ip string) (n string, ok bool) {
	n, ok = nameTable[ip]

	return
}
