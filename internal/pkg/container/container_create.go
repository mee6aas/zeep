package container

import (
	"context"

	"github.com/mee6aas/zeep/api"

	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

// CreateCreatedBody OK response to ContainerCreate operation
type CreateCreatedBody struct {
	ID string // ID of container
}

// Create creates a container based in the given configuration.
func Create(ctx context.Context, config *Config) (response CreateCreatedBody, err error) {
	var (
		res cont.ContainerCreateCreatedBody
	)

	if res, err = engineClient.ContainerCreate(ctx, &cont.Config{
		Image: config.Image,
		Env:   []string{},
	}, &cont.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: config.Storage,
				Target: api.KyleStorage,
				BindOptions: &mount.BindOptions{
					Propagation: mount.PropagationShared,
				},
			},
		},
	}, nil, ""); err != nil {
		return
	}

	response.ID = res.ID

	return
}
