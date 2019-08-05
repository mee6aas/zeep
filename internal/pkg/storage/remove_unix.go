// +build linux freebsd openbsd darwin

package storage

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Remove removes a storage.
func (s Storage) Remove() (e error) {
	var (
		errs []string
	)

	if err := unix.Unmount(s.path, unix.MNT_DETACH); err != nil &&
		err.Error() != "invalid argument" &&
		err.Error() != "no such file or directory" {
		err = errors.Wrap(err, "Failed to unmount")
		errs = append(errs, err.Error())
	}

	if err := os.Remove(s.path); err != nil && !os.IsNotExist(err) {
		err = errors.Wrap(err, "Failed to remove")
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		e = fmt.Errorf(strings.Join(errs, "\n"))
	}

	return
}

// RemoveDetach removes a storage and warns if failed not return error.
func (s Storage) RemoveDetach() {
	if err := s.Remove(); err != nil {
		// warn
	}
}
