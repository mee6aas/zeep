package acts

import (
	"github.com/mee6aas/zeep/pkg/activity"
)

// Read returns activity in collection.
func Read(username string, actName string) (a activity.Activity, ok bool) {
	var (
		as map[string]activity.Activity
	)

	if as, ok = activities[username]; !ok {
		return
	}

	a, ok = as[actName]

	return
}
