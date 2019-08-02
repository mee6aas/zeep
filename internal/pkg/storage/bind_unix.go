// +build linux freebsd openbsd darwin

package storage

import (
	"golang.org/x/sys/unix"
)

const (
	// MSReadOnly equals unix.MS_RDONLY
	MSReadOnly = unix.MS_RDONLY
)

// Bind mounts `src` direcotry to `trg`.
func Bind(trg string, src string, flag uintptr) (s Storage, e error) {
	f := flag | unix.MS_BIND

	if e = unix.Mount(src, trg, "", f, ""); e != nil {
		return
	}

	s.path = trg

	return
}
