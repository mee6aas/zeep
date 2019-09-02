package container

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
)

// Start starts this container.
func (c *Container) Start(ctx context.Context) (e error) {
	if e = engineClient.ContainerStart(ctx, c.id, types.ContainerStartOptions{}); e != nil {
		return
	}

	{
		info, e := engineClient.ContainerInspect(ctx, c.id)
		if e != nil {
			return e
		}

		es, ok := info.NetworkSettings.Networks[agentNet]
		if !ok {
			e = errors.New("Container has no network settings")
			return e
		}

		c.ip = es.IPAddress
	}

	return
}
