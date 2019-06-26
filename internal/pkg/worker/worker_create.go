package worker

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/container"
	"github.com/mee6aas/zeep/internal/pkg/storage"

	"github.com/pkg/errors"
)

// CreateCreatedBody OK reponse to ContainerCreate operation
type CreateCreatedBody struct {
	ID string // ID of container
}

// Create creates a worker based in the given configuration.
func Create(ctx context.Context, config *Config) (response CreateCreatedBody, err error) {
	var (
		strRes storage.CreateCreatedBody
		cntRes container.CreateCreatedBody
	)

	if strRes, err = storage.Create(storage.CreateConfig{Size: config.Size}); err != nil {
		err = errors.Wrap(err, "Failed to create storage")
		return
	}

	if cntRes, err = container.Create(ctx, &container.Config{
		Image:   config.Image,
		Storage: strRes.Path,
	}); err != nil {
		err = errors.Wrap(err, "Failed to create container")
		return
	}

	config.path = strRes.Path

	workers[cntRes.ID] = *config

	response.ID = cntRes.ID

	return
}
