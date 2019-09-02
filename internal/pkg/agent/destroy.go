package agent

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
)

// Destroy removes agent resources.
func Destroy(ctx context.Context) (e error) {
	if !isSetup {
		return
	}

	// TODO: currently only last error is preserved.

	// TODO: stop grpc gracfully after disconnect allocs before remove

	log.Debug("Destorying acts")
	e = acts.Destroy(ctx)
	if e == nil {
		log.Debug("Acts destroyed")
	} else {
		log.WithError(e).Warn("Failed to destroy acts")
	}

	log.Debug("Destorying allocs")
	e = allocs.Destroy(ctx)
	if e == nil {
		log.Debug("Allocs destroyed")
	} else {
		log.WithError(e).Warn("Failed to destroy allocs")
	}

	log.Debug("Destorying worker pool")
	workerPool.Destroy(ctx)
	if e == nil {
		log.Debug("Worker pool destroyed")
	} else {
		log.WithError(e).Warn("Failed to destory worker pool")
	}

	log.Debug("Stopping gRPC server")
	gRPCServer.Stop()
	log.Debug("gRPC server stop")

	isSetup = false

	return
}
