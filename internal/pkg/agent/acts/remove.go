package acts

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mee6aas/zeep/pkg/activity"
)

// Remove deletes activity from collection with given username and actName.
func Remove(username string, actName string) (e error) {
	var (
		ok   bool
		acts map[string]activity.Activity
	)

	if acts, ok = activities[username]; !ok {
		e = errors.New("Not found")
		return
	}

	if e = os.RemoveAll(filepath.Join(rootDirPath, actName)); e != nil {
		return
	}

	delete(acts, actName)

	return
}
