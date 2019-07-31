package acts

import (
	"os"
	"path/filepath"

	"github.com/mee6aas/zeep/pkg/activity"
)

const (
	// DefaultRootDirPath indicates the path of directory that used by acts.
	DefaultRootDirPath = "/usr/zeep/acts"
)

var (
	isSetup     bool
	rootDirPath string

	//                username    actID
	entries = make(map[string]map[string]activity.Activity)
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
	// RootDirPath indicates the path of directory that used by acts.
	RootDirPath string
}

func normalizeConfig(config Config) (c Config) {
	c = config
	c.RootDirPath = filepath.Clean(c.RootDirPath)
	if c.RootDirPath == "" {
		c.RootDirPath = DefaultRootDirPath
	}

	return
}

// Setup initializes acts root dir and initial activities.
func Setup(config Config) (e error) {
	if IsSetup() {
		Destroy()
	}

	config = normalizeConfig(config)
	rootDirPath = config.RootDirPath

	if _, e = os.Stat(rootDirPath); os.IsNotExist(e) {
		os.MkdirAll(rootDirPath, os.ModePerm)
	} else if e != nil {
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

	rootDirPath = ""
	entries = make(map[string]map[string]activity.Activity)

	isSetup = false

	return
}
