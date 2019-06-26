// +build linux freebsd openbsd darwin

package storage

import (
	"os"

	"golang.org/x/sys/unix"
)

// Remove removes a storage.
func Remove(path string) (err error) {
	if err = unix.Unmount(path, unix.MNT_DETACH); err != nil {
		return
	}

	if err = os.Remove(path); err != nil {
		return
	}

	return
}
