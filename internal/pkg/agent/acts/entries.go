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
func EntriesInUsername(username string) (as map[string]activity.Activity) {
	e, ok := activities[username]
	if !ok {
		return
	}

	as = e

	return
}
