package acts

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mee6aas/zeep/pkg/activity"
)

// Remove deletes activity from collection with specified username and ID.
func Remove(username string, actID string) (e error) {
	var (
		ok   bool
		acts map[string]activity.Activity
	)

	if acts, ok = entries[username]; !ok {
		e = errors.New("Not found")
		return
	}

	if e = os.RemoveAll(filepath.Join(rootDirPath, actID)); e != nil {
		return
	}

	delete(acts, actID)

	return
}
