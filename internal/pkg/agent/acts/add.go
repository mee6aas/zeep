package acts

import (
	"errors"
	"path/filepath"

	"github.com/otiai10/copy"

	"github.com/mee6aas/zeep/pkg/activity"
)

// Add inserts activity to collection with given username and id.
func Add(username string, actName string, actDirPath string) (e error) {
	if !IsSetup() {
		e = errors.New("Acts not setup")
		return
	}

	var (
		ok   bool
		act  activity.Activity
		acts map[string]activity.Activity
	)

	if act, e = activity.UnmarshalFromDir(actDirPath); e != nil {
		return
	}

	if e = copy.Copy(actDirPath, filepath.Join(rootDirPath, username, actName)); e != nil {
		return
	}

	if acts, ok = activities[username]; !ok {
		acts = make(map[string]activity.Activity)
	}

	acts[actName] = act
	activities[username] = acts

	return
}
