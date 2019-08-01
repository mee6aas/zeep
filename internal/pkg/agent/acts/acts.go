package acts

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/pkg/activity"
)

var (
	isSetup     bool
	rootDirPath string

	//                   username    actID
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

// Setup initializes acts root dir and initial activities.
func Setup(config Config) (e error) {
	if IsSetup() {
		Destroy()
	}

	if rootDirPath, e = ioutil.TempDir("", ""); e != nil {
		e = errors.Wrap(e, "Failed to create root directory")
		return
	}

	isSetup = true

	return
}

// Destroy removes acts resources.
func Destroy() (e error) {
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
