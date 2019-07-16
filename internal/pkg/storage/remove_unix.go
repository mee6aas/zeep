// +build linux freebsd openbsd darwin

package storage

import (
	"os"

	"golang.org/x/sys/unix"
)

// Remove removes a storage.
func (s *Storage) Remove() (err error) {
	if err = unix.Unmount(s.path, unix.MNT_DETACH); err != nil {
		return
	}

	if err = os.Remove(s.path); err != nil {
		return
	}

	return
}
