package acts

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/mholt/archiver"
	"github.com/otiai10/copy"

	"github.com/mee6aas/zeep/pkg/activity"
)

// AddFromDir adds an activity at the specified directory with the specified username and activity name to collection.
func AddFromDir(username string, actName string, actDirPath string) (e error) {
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
		e = errors.Wrapf(e, "Failed to unmarshal from %s", actDirPath)
		return
	}

	src := filepath.Join(rootDirPath, username, actName)
	if e = copy.Copy(actDirPath, src); e != nil {
		e = errors.Wrapf(e, "Failed to copy from %s to %s", src, actDirPath)
		return
	}

	if acts, ok = activities[username]; !ok {
		acts = make(map[string]activity.Activity)
	}

	act.Owner = username
	act.Name = actName
	act.AddedDate = time.Now().String()

	acts[actName] = act
	activities[username] = acts

	return
}

// AddFromTarGz adds an activity at the specified Gzip with the specified username and activity name to collection.
// The file `actTarGzPath` must be gzipped tarball.
func AddFromTarGz(username string, actName string, actTarGzPath string) (e error) {
	var (
		trg string
	)

	if trg, e = ioutil.TempDir("", ""); e != nil {
		e = errors.Wrap(e, "Failed to create temporal directory")
		return
	}
	defer os.RemoveAll(trg)

	tgz := archiver.NewTarGz()
	tgz.ImplicitTopLevelFolder = true
	if e = tgz.Unarchive(actTarGzPath, trg); e != nil {
		e = errors.Wrapf(e, "Failed to unarchive %s", actTarGzPath)
		return
	}

	if e = AddFromDir(username, actName, trg); e != nil {
		return
	}

	return
}

// AddFromHTTP adds an activity at the specified HTTP address with specified username and activity name to collection.
// The address `actAddr` must be http address and must response gzipped tarball as the attachment.
func AddFromHTTP(username string, actName string, actAddr string) (e error) {
	var (
		res *http.Response
		trg *os.File
	)

	if res, e = http.Get(actAddr); e != nil {
		e = errors.Wrapf(e, "Failed to GET %s", actAddr)
		return
	}
	defer res.Body.Close()

	if trg, e = ioutil.TempFile("", ""); e != nil {
		e = errors.Wrap(e, "Failed to create temporal directory")
		return
	}
	defer os.Remove(trg.Name())

	if _, e = io.Copy(trg, res.Body); e != nil {
		e = errors.Wrapf(e, "Failed to write response from %s", actAddr)
		return
	}

	if e = AddFromTarGz(username, actName, trg.Name()); e != nil {
		return
	}

	return
}
