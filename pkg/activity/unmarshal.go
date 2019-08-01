package activity

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Unmarshal parses activity configuration.
func Unmarshal(act []byte) (a Activity, e error) {
	e = json.Unmarshal(act, &a)
	a = Normalize(a)

	return
}

// UnmarshalFromFile parses activity manifest from given file.
func UnmarshalFromFile(actDesc string) (a Activity, e error) {
	var (
		act []byte
	)

	if act, e = ioutil.ReadFile(actDesc); e != nil {
		return
	}

	a, e = Unmarshal(act)

	return
}

// UnmarshalFromDir parses activity manifest in directory with default manifest name.
func UnmarshalFromDir(actDir string) (a Activity, e error) {
	var (
		fs []os.FileInfo
	)

	if fs, e = ioutil.ReadDir(actDir); e != nil {
		return
	}

	for _, f := range fs {
		if f.Name() == DefaultActivityManifestName {
			a, e = UnmarshalFromFile(filepath.Join(actDir, f.Name()))

			return
		}
	}

	e = errors.New("Not found")

	return
}
