package v1

import (
	"context"

	acts "github.com/mee6aas/zeep/internal/pkg/var/activities"
)

// AddRequested is invoked when the invoker requests to add an activity.
func (h Handle) AddRequested(
	_ context.Context,
	username string,
	actName string,
	actDirPath string,
) (e error) {
	e = acts.AddFromDir(username, actName, actDirPath)

	return
}
