package acts

import "github.com/mee6aas/zeep/pkg/activity"

// Entries returns all the activities in collection.
func Entries() (as map[string]activity.Activity) {
	for _, e := range activities {
		for k, v := range e {
			as[k] = v
		}
	}

	return
}

// EntriesInUsername returns the activities in given username.
func EntriesInUsername(username string) (as []activity.Activity, ok bool) {
	var (
		es map[string]activity.Activity
	)

	if es, ok = activities[username]; !ok {
		return
	}

	as = make([]activity.Activity, 0, len(es))

	for _, v := range es {
		as = append(as, v)
	}

	return
}
