package storage

var (
	storageRoot = "/tmp/"
)

// Storage describes storage.
type Storage struct {
	path string // Path of this storage
	size uint64 // Size of this storage in KiB
}

// Path returns the path of this storage.
func (s *Storage) Path() string { return s.path }

// Config holds the configuration for the storage.
type Config struct {
	Size uint64 // Size of sotrage to create in KiB.
}
