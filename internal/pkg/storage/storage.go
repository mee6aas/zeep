package storage

import (
	"os"
	"path/filepath"

	"github.com/mee6aas/zeep/api"
)

var (
	rootPathOnHost string
)

// Storage describes storage.
type Storage struct {
	path string // Path of this storage
	size uint64 // Size of this storage in KiB
}

// Path returns the path of this storage.
func (s Storage) Path() string { return s.path }

// PathOnHost returns the path of this storage on host.
func (s Storage) PathOnHost() string {
	rel, _ := filepath.Rel(os.TempDir(), s.path)
	return filepath.Join(rootPathOnHost, rel)
}

// Config holds the configuration for the storage.
type Config struct {
	Size uint64 // Size of sotrage to create in KiB.
}

func init() {
	if rootPathOnHost = os.Getenv(api.AgentTmpDirPathEnvKey); rootPathOnHost == "" {
		rootPathOnHost = os.TempDir()
	}
}
