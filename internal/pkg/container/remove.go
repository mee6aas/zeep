package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

// Remove removes a container.
func (c Container) Remove(ctx context.Context) (err error) {
	if err = engineClient.ContainerRemove(ctx, c.id, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		errors.Wrapf(err, "Failed to remove container %s", c.id)
		return
	}

	return
}
