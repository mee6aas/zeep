package worker

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/internal/pkg/container"
	"github.com/mee6aas/zeep/internal/pkg/storage"
)

// Worker describes worker.
type Worker struct {
	InvokeeVersion string // Version of invokee service that assigned to this worker

	container    container.Container // Container descriptor
	storage      storage.Storage     // Storage descriptor
	isAllocated  bool                // Is taskAssigner set
	taskAssigner TaskAssigner        // Pass task to connected stream
	isAssigned   bool                // Is task assigned
}

// ID returns the ID of the container used by this worker.
func (w Worker) ID() string { return w.container.ID() }

// Container returns a descriptor for the container used by this worker.
func (w Worker) Container() container.Container { return w.container }

// Storage returns a descriptor for the storage used by this worker
func (w Worker) Storage() storage.Storage { return w.storage }

// IsAllocated checks if the task is allocated to this worker.
func (w Worker) IsAllocated() bool { return w.isAllocated }

// IsAssigned checks if the task is assigned to this worker.
func (w Worker) IsAssigned() bool { return w.isAssigned }

// Resolve set isAssigned flag to false
func (w *Worker) Resolve() { w.isAssigned = false }

// Dealloc set isAllocated flag to false
func (w *Worker) Dealloc() { w.isAllocated = false }

// Config holds the configuration for the worker.
type Config struct {
	Image string // Image to use
	Size  uint64 // Size limit of the storage to mount
}

// NewWorker creates a new worker based on the given configuration
// and returns its descriptor.
func NewWorker(ctx context.Context, config Config) (worker Worker, e error) {
	var (
		sto  storage.Storage
		cont container.Container
	)

	if sto, e = storage.NewStorage(storage.Config{Size: config.Size}); e != nil {
		e = errors.Wrap(e, "Failed to create storage")
		return
	}

	os.Mkdir(filepath.Join(sto.Path(), filepath.Base(api.ActivityResource)), 755)
	// os.Mkdir(filepath.Join(sto.Path(), filepath.Base(api.WorkflowStorage)), 755)

	if cont, e = container.NewContainer(ctx, container.Config{
		Image:   config.Image,
		Storage: sto.Path(),
	}); e != nil {
		sto.Remove()
		e = errors.Wrap(e, "Failed to create container")
		return
	}

	worker.container = cont
	worker.storage = sto

	return
}
