package acts_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/agent/acts"
)

func TestSetUpAndDestory(t *testing.T) {
	var (
		err error
	)

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.Destroy(context.Background()); err != nil {
		t.Fatalf("Failed to destory acts: %v", err)
	}
}

func TestAddAndRemove(t *testing.T) {
	var (
		err error

		username = "Jerry"
		actName  = "golf"
	)

	if err = acts.AddFromDir(username, actName, "./testdata/valid"); err == nil {
		t.Fatalf("Expected to fail to add activity")
	}

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	defer func() {
		if e := acts.Destroy(context.Background()); e != nil {
			t.Logf("Failed to destroy acts: %v", e)
		}
	}()

	if err = acts.AddFromDir(username, actName, "./testdata/valid"); err != nil {
		t.Fatalf("Failed to add activity")
	}

	if _, ok := acts.Read(username, actName); !ok {
		t.Fatalf("Expected that the activity %s/%s is exists", username, actName)
	}

	if _, err = os.Stat(filepath.Join(acts.RootDirPath(), username, actName)); os.IsNotExist(err) {
		t.Fatalf("Expected that the activity is added at %s", filepath.Join(acts.RootDirPath(), username, actName))
	}

	if err = acts.Remove(username, actName); err != nil {
		t.Fatalf("Failed to remove activity")
	}

	if err = acts.Destroy(context.Background()); err != nil {
		t.Fatalf("Failed to destroy acts")
	}
}

func TestAddFromTarGz(t *testing.T) {
	var (
		err error

		username = "Summer"
		actName  = "fame"
	)

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.AddFromTarGz(username, actName, "./testdata/valid/valid.tar.gz"); err != nil {
		t.Fatalf("Failed to add gzipped tarball activity: %v", err)
	}

	if _, ok := acts.Read(username, actName); !ok {
		t.Fatalf("Expected that the activity %s/%s is exists", username, actName)
	}

	if _, err = os.Stat(filepath.Join(acts.RootDirPath(), username, actName)); os.IsNotExist(err) {
		t.Fatalf("Expected that the activity is added at %s", filepath.Join(acts.RootDirPath(), username, actName))
	}

	if err = acts.Remove(username, actName); err != nil {
		t.Fatalf("Failed to remove activity")
	}

	if err = acts.Destroy(context.Background()); err != nil {
		t.Fatalf("Failed to destroy acts")
	}
}

func TestAddFromHTTP(t *testing.T) {
	var (
		err error
		srv = &http.Server{Addr: ":5125"}

		username = "beth"
		actName  = "individual"
	)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var (
			err  error
			src  *os.File
			size string
		)

		if src, err = os.Open("./testdata/valid/valid.tar.gz"); err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}

		if info, e := src.Stat(); e == nil {
			size = strconv.FormatInt(info.Size(), 10)
		} else {
			t.Fatalf("Failed to stat test file: %v", e)
		}

		w.Header().Set("Content-Disposition", "attachment; filename=valid.tar.gz")
		w.Header().Set("Content-Type", "application/gzip")
		w.Header().Set("Content-Length", size)

		if _, err = io.Copy(w, src); err != nil {
			t.Fatalf("Failed to response file: %v", err)
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Fatalf("Failed to serve: %v", err)
		}
	}()

	if err = acts.Setup(acts.Config{}); err != nil {
		t.Fatalf("Failed to setup acts: %v", err)
	}

	if err = acts.AddFromHTTP(username, actName, "http://localhost:5125/"); err != nil {
		t.Fatalf("Failed to add gzipped tarball activity: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	if err = srv.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}
	cancel()

	if _, ok := acts.Read(username, actName); !ok {
		t.Fatalf("Expected that the activity %s/%s is exists", username, actName)
	}

	if _, err = os.Stat(filepath.Join(acts.RootDirPath(), username, actName)); os.IsNotExist(err) {
		t.Fatalf("Expected that the activity is added at %s", filepath.Join(acts.RootDirPath(), username, actName))
	}

	if err = acts.Remove(username, actName); err != nil {
		t.Fatalf("Failed to remove activity")
	}

	if err = acts.Destroy(context.Background()); err != nil {
		t.Fatalf("Failed to destroy acts")
	}
}
