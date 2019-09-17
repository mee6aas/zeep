package v1

import (
	"context"

	acts "github.com/mee6aas/zeep/internal/pkg/var/activities"
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
