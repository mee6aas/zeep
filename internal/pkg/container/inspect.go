package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Inspect returns detailed information on this container.
func (c Container) Inspect(ctx context.Context) (i types.ContainerJSON, e error) {
	i, e = engineClient.ContainerInspect(ctx, c.id)
	return
}
