package v1

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
)

// RegisterRequested is invoked when the invoker requests an activity registration.
func (h Handle) RegisterRequested(
	_ context.Context,
	username string,
	actName string,
	actDirPath string,
) (e error) {
	e = acts.AddFromDir(username, actName, actDirPath)

	return
}
