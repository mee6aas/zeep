package acts

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

// Setup initializes the collection.
func Setup(config Config) (e error) {
	if IsSetup() {
		return
	}

	if rootDirPath, e = ioutil.TempDir("", "zeep-act-"); e != nil {
		e = errors.Wrap(e, "Failed to create root directory")
		return
	}

	isSetup = true

	return
}
