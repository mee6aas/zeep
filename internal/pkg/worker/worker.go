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

	image        string
	container    container.Container // Container descriptor
	storage      storage.Storage     // Storage descriptor
	isAllocated  bool                // Is taskAssigner set
	taskAssigner TaskAssigner        // Pass task to connected stream
	isAssigned   bool                // Is task assigned
}

// ID returns the ID of the container used by this worker.
func (w Worker) ID() string { return w.container.ID() }

// IP returns the IP of the container used by this worker.
func (w Worker) IP() string { return w.container.IP() }

// Image returns the image of the container used by this worker.
func (w Worker) Image() string { return w.image }

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

// Config holds the configuration for the worker.
type Config struct {
	Image string // Image to use
	Size  uint64 // Size limit of the storage to mount
}

// NewWorker creates a new worker based on the given configuration
// and returns its descriptor.
func NewWorker(ctx context.Context, conf Config) (worker Worker, e error) {
	var (
		sto  storage.Storage
		cont container.Container
	)

	if sto, e = storage.NewStorage(storage.Config{Size: conf.Size}); e != nil {
		e = errors.Wrap(e, "Failed to create storage")
		return
	}
	defer func() {
		if e != nil {
			sto.RemoveDetach()
		}
	}()

	os.Mkdir(filepath.Join(sto.Path(), filepath.Base(api.ActivityResource)), 755)
	// os.Mkdir(filepath.Join(sto.Path(), filepath.Base(api.WorkflowStorage)), 755)

	if cont, e = container.NewContainer(ctx, container.Config{
		Image:   conf.Image,
		Storage: sto.PathOnHost(),
	}); e != nil {
		e = errors.Wrap(e, "Failed to create container")
		return
	}

	worker.image = conf.Image
	worker.container = cont
	worker.storage = sto

	return
}
