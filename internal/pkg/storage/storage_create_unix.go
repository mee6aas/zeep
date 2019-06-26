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

// Create creates a new storage based in the given configuration.
func Create(config CreateConfig) (response CreateCreatedBody, err error) {
	var (
		trg string // path of the storage to create
		opt string // option of the tmpfs
	)

	trg = filepath.Join(storageRoot, uuid.New().String())

	if err = os.Mkdir(trg, os.ModePerm); err != nil {
		return
	}

	opt = fmt.Sprintf("size=%d", config.Size)

	if err = unix.Mount("tmpfs", trg, "tmpfs", 0, opt); err != nil {
		err = errors.Wrap(err, "failed to mount tmpfs")
		if e := os.Remove(trg); e != nil {
			e = errors.Wrapf(e, "failed to remove mountpoint %s", trg)
			err = errors.Wrap(err, e.Error())
		}
		return
	}

	response.Path = trg

	return
}
