package activity

import (
	"strings"

	"github.com/mee6aas/zeep/api"
)

// Normalize returns normalized activity based on given activity.
func Normalize(act Activity) (a Activity) {
	a = act

	if !strings.Contains(a.Runtime, "/") {
		a.Runtime = api.Mee6aaSDockerOrgName + "/" + a.Runtime
	}

	return
}
