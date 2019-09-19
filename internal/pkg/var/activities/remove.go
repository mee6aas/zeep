package acts

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mee6aas/zeep/pkg/activity"
)

// Remove remvoes an activity with the specified username and activity name from the collection.
func Remove(username string, actName string) (e error) {
	var (
		ok   bool
		acts map[string]activity.Activity
	)

	if acts, ok = activities[username]; !ok {
		e = errors.New("Not found")
		return
	}

	if _, ok = acts[actName]; !ok {
		e = errors.New("Not found")
		return
	}

	if e = os.RemoveAll(filepath.Join(rootDirPath, actName)); e != nil {
		if !os.IsNotExist(e) {
			return
		}
	}

	delete(acts, actName)

	return
}
