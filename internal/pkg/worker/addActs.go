package worker

import (
	"path/filepath"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/storage"
)

// AddActs binds given directory path to the worker resource directory.
func (w Worker) AddActs(actDirPath string) (e error) {
	// it will be unmounted by manually in Remove
	_, e = storage.Bind(filepath.Join(w.storage.Path(), filepath.Base(api.ActivityResource)), actDirPath, 0)
	return
}
