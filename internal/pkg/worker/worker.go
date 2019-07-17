package worker

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/container"
	"github.com/mee6aas/zeep/internal/pkg/storage"
	"github.com/pkg/errors"
)

// TaskAssigner passes task to its task operator.
type TaskAssigner interface {
	Assign(context.Context, interface{}) error
}

// Worker describes worker.
type Worker struct {
	InvokeeVersion string // Version of invokee service that assigned to this worker

	container    container.Container // Container descriptor
	storage      storage.Storage     // Storage descriptor
	isConnected  bool                // Is taskAssigner set
	taskAssigner TaskAssigner        // Pass task to connected stream
}

// ID returns the ID of the container used by this worker.
func (w Worker) ID() string { return w.container.ID() }

// Container returns a descriptor for the container used by this worker.
func (w Worker) Container() container.Container { return w.container }

// Storage returns a descriptor for the storage used by this worker
func (w Worker) Storage() storage.Storage { return w.storage }

// IsConnected checks if this worker connted to its task operator.
func (w Worker) IsConnected() bool { return w.isConnected }

// Connect connects task assigner to this worker.
// By this behavior, this worker is connected to the task operator.
func (w Worker) Connect(ta TaskAssigner) (err error) {
	if w.IsConnected() {
		err = errors.New("Already connected worker")
		return
	}

	w.isConnected = true
	w.taskAssigner = ta

	return
}

// Config holds the configuration for the worker.
type Config struct {
	Image string // Image to use
	Size  uint64 // Size limit of the storage to mount
}

// NewWorker creates a new worker based on the given configuration
// and returns its descriptor.
func NewWorker(ctx context.Context, config Config) (worker Worker, err error) {
	var (
		stor storage.Storage
		cont container.Container
	)

	if stor, err = storage.NewStorage(storage.Config{Size: config.Size}); err != nil {
		err = errors.Wrap(err, "Failed to create storage")
		return
	}

	if cont, err = container.NewContainer(ctx, container.Config{
		Image:   config.Image,
		Storage: stor.Path(),
	}); err != nil {
		stor.Remove()
		err = errors.Wrap(err, "Failed to create container")
		return
	}

	worker.container = cont
	worker.storage = stor

	return
}
