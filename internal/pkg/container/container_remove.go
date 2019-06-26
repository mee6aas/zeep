package container

import (
	"context"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
)

// Remove removes a container.
func Remove(ctx context.Context, containerID string) (err error) {
	if err = engineClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		errors.Wrapf(err, "Failed to remove container %s", containerID)
		return
	}

	return
}
