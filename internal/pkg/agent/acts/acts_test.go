package acts_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
)

func TestSetUpAndDestory(t *testing.T) {
	var (
		err error
	)

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.Destroy(); err != nil {
		t.Fatalf("Failed to destory acts: %v", err)
	}
}

func TestAddAndRemove(t *testing.T) {
	var (
		err error

		username = "Jerry"
		actName  = "golf"
	)

	if err = acts.Add(username, actName, "./testdata/valid"); err == nil {
		t.Fatalf("Expected to fail to add activity")
	}

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	defer func() {
		if e := acts.Destroy(); e != nil {
			t.Logf("Failed to destroy acts: %v", e)
		}
	}()

	if err = acts.Add(username, actName, "./testdata/valid"); err != nil {
		t.Fatalf("Failed to add activity")
	}

	if _, err = os.Stat(filepath.Join(acts.RootDirPath(), username, actName)); os.IsNotExist(err) {
		t.Fatalf("Expected that the activity is added at %s", filepath.Join(acts.RootDirPath(), username, actName))
	}

	if err = acts.Remove(username, actName); err != nil {
		t.Fatalf("Failed to remove activity")
	}
}
