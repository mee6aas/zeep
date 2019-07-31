package acts_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
)

func TestSetUpAndDestory(t *testing.T) {
	var (
		err error
		d   string
	)

	if d, err = ioutil.TempDir("", ""); err != nil {
		t.Fatalf("Failed to create temporal directory: %v", err)
	}

	t.Logf("Temporal direcotry %s created", d)

	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Logf("Failed to remove temporal direcotry %s: %v", d, err)
		}
	}()

	if err = acts.Setup(acts.Config{
		RootDirPath: d,
	}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.Destroy(); err != nil {
		t.Fatalf("Failed to destory acts: %v", err)
	}
}

func TestAddAndRemove(t *testing.T) {
	var (
		err error
		d   string

		username = "Jerry"
		actID    = "golf"
	)

	if err = acts.Add(username, actID, "./testdata/valid"); err == nil {
		t.Fatalf("Expected to fail to add activity")
	}

	if d, err = ioutil.TempDir("", ""); err != nil {
		t.Fatalf("Failed to create temporal directory: %v", err)
	}

	t.Logf("Temporal direcotry %s created", d)

	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Logf("Failed to remove temporal direcotry %s: %v", d, err)
		}
	}()

	if err = acts.Setup(acts.Config{
		RootDirPath: d,
	}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.Add(username, actID, "./testdata/valid"); err != nil {
		t.Fatalf("Failed to add activity")
	}

	if _, err = os.Stat(filepath.Join(d, actID)); os.IsNotExist(err) {
		t.Fatalf("Expected that the activity is added at %s", filepath.Join(d, actID))
	}

	if err = acts.Remove(username, actID); err != nil {
		t.Fatalf("Failed to remove activity")
	}
}
