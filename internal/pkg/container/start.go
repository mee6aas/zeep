package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/api"
)

// Start starts this container.
func (c *Container) Start(ctx context.Context) (err error) {
	var (
		ok   bool
		info types.ContainerJSON
		es   *network.EndpointSettings
	)

	if err = engineClient.ContainerStart(ctx, c.id, types.ContainerStartOptions{}); err != nil {
		return
	}

	if info, err = engineClient.ContainerInspect(ctx, c.id); err != nil {
		return
	}

	if es, ok = info.NetworkSettings.Networks[api.NetworkName]; !ok {
		err = errors.New("Container has no network settings")
		return
	}

	c.ip = es.IPAddress

	return
}
