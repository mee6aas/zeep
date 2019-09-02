package v1

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
)

// RemoveRequested is invoked when the invoker requests to remove an activity.
func (h Handle) RemoveRequested(
	_ context.Context,
	username string,
	actName string,
) (e error) {
	// TODO: do not affect to running activity

	return acts.Remove(username, actName)
}
