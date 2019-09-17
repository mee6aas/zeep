package acts

import (
	"context"
	"os"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/pkg/activity"
)

// Destroy removes acts resources.
func Destroy(_ context.Context) (e error) {
	if !IsSetup() {
		return
	}

	if e = os.RemoveAll(rootDirPath); e != nil {
		e = errors.Wrapf(e, "Failed to remove root directory")
		return
	}

	rootDirPath = ""
	activities = make(map[string]map[string]activity.Activity)

	isSetup = false

	return
}
