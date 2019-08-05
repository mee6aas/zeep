// +build linux freebsd openbsd darwin

package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// NewStorage creates a new storage based on the given configuration
// and returns its descriptor.
func NewStorage(config Config) (s Storage, e error) {
	var (
		trg string // path of the storage to create
		opt string // option for tmpfs
	)

	trg = filepath.Join(storageRoot, uuid.New().String())

	if e = os.Mkdir(trg, os.ModePerm); e != nil {
		return
	}
	defer func() {
		if e != nil {
			os.Remove(trg)
		}
	}()

	opt = fmt.Sprintf("size=%d", config.Size)

	if config.Size == 0 {
		opt = ""
	}

	if e = unix.Mount("tmpfs", trg, "tmpfs", 0, opt); e != nil {
		e = errors.Wrap(e, "failed to mount tmpfs")
		return
	}

	s.path = trg
	s.size = config.Size

	return
}
