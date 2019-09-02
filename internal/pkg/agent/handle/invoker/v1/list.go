package v1

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/pkg/activity"
)

// ListRequested is invoked when the invoker requests to list the activities.
func (h Handle) ListRequested(
	_ context.Context,
	username string,
) (out []activity.Activity, e error) {
	var (
		ok bool
	)

	out, ok = acts.EntriesInUsername(username)

	if !ok {
		e = errors.New("Not found")
		return
	}

	return
}
