package handover

var (
	//			 IP		not used
	reserved map[string]bool
)

// IsReserved checks if the specified address is being handovered.
func IsReserved(addr string) (ok bool) {
	_, ok = reserved[addr]

	return
}

// Reserve marks that the specified address is being handovered.
func Reserve(addr string) bool {
	if IsReserved(addr) {
		return false
	}

	reserved[addr] = true

	return true
}

// Resolve unmarks that the specified address.
func Resolve(addr string) bool {
	if !IsReserved(addr) {
		return false
	}

	delete(reserved, addr)

	return true
}
