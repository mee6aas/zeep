package worker

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/storage"
)

// AddActs binds given directory path to the worker resource directory.
func (w Worker) AddActs(actDirPath string) (e error) {
	// it will be unmounted by manually in worker.Remove
	_, e = storage.Bind(filepath.Join(w.storage.Path(), filepath.Base(api.ActivityResource)), actDirPath, 0)

	log.WithFields(log.Fields{
		"ID":   w.ID(),
		"what": actDirPath,
	}).Debug("Activity bound")

	return
}
