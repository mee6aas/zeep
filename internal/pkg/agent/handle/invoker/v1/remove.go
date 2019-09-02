package v1

import (
	"context"
)

// RemoveRequested is invoked when the invoker requests to remove an activity.
func (h Handle) RemoveRequested(
	_ context.Context,
	username string,
	actName string,
) (e error) {

	return
}
