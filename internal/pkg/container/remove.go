package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Remove removes a container.
func (c Container) Remove(ctx context.Context) (e error) {
	if e = engineClient.ContainerRemove(ctx, c.id, types.ContainerRemoveOptions{
		Force: true,
	}); e != nil {
		return
	}

	return
}
