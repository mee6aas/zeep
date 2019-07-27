package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Inspect returns detailed information on this container.
func (c Container) Inspect(ctx context.Context) (info types.ContainerJSON, err error) {
	info, err = engineClient.ContainerInspect(ctx, c.id)
	return
}
