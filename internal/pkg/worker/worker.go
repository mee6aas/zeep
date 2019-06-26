package worker

var (
	workers map[string]Config
)

// Config holds the configuration for the worker.
type Config struct {
	Image string // Name of the image worker use
	Size  uint32 // Size limit of the storage to create

	path string // Path of the storage
}

// TODO: resource reuse policy
