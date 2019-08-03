package v1

import (
	"context"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
	"github.com/mee6aas/zeep/pkg/activity"
)

func (h Handle) RegisterRequested(
	_ context.Context,
	username string,
	actDirPath string,
) (e error) {
	var (
		a activity.Activity
	)

	a, e = activity.UnmarshalFromDir(actDirPath)
	e = acts.Add(username, a.Name, actDirPath)

	return
}
