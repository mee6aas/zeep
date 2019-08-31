package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
)

// Start starts this container.
func (c *Container) Start(ctx context.Context) (e error) {
	var (
		ok   bool
		info types.ContainerJSON
		es   *network.EndpointSettings
	)

	if e = engineClient.ContainerStart(ctx, c.id, types.ContainerStartOptions{}); e != nil {
		return
	}

	if info, e = engineClient.ContainerInspect(ctx, c.id); e != nil {
		return
	}

	if es, ok = info.NetworkSettings.Networks[agentNet]; !ok {
		e = errors.New("Container has no network settings")
		return
	}

	c.ip = es.IPAddress

	return
}
