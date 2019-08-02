// +build linux freebsd openbsd darwin

package storage

import (
	"golang.org/x/sys/unix"
)

// Unmount unmounts given directory.
func Unmount(trg string) (e error) {
	if e = unix.Unmount(trg, unix.MNT_DETACH); e != nil {
		return
	}

	return
}
