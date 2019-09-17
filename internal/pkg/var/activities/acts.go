package acts

import (
	"github.com/mee6aas/zeep/pkg/activity"
)

var (
	isSetup     bool
	rootDirPath string

	// activities holds activity manifest.
	//                   username   actName
	activities = make(map[string]map[string]activity.Activity)
)

// TODO:
// - limit activity resource size limit
// - limit activity rootfolder quota limit
// - remove longtime unused activity
// - add longtime unused activity eviction(remove) strategy considering quota limit.

// IsSetup check if acts is setup.
func IsSetup() bool { return isSetup }

// RootDirPath returns the path of directory that used by acts.
func RootDirPath() string { return rootDirPath }

// Config holds the configuration for the acts.
type Config struct {
}
