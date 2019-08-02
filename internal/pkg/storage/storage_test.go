package storage_test

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/mee6aas/zeep/internal/pkg/storage"
)

var (
	testStorageCreateFailed = false

	testStorage      Storage
	testStorageQuota = uint64(1024 * 4) // 4KiB
)

func TestStorageCreate(t *testing.T) {
	var (
		err error
	)

	if testStorage, err = NewStorage(Config{
		Size: testStorageQuota * 2,
	}); err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
}

func TestStorageQuota(t *testing.T) {
	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	var (
		err error
		trg *os.File
	)

	if trg, err = ioutil.TempFile(testStorage.Path(), ""); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := io.CopyN(trg, rand.Reader, int64(testStorageQuota)); err != nil {
		t.Fatalf("expected to create %d size file: %v", testStorageQuota, err)
	}

	if _, err := io.CopyN(trg, rand.Reader, int64(testStorageQuota*3)); err == nil {
		t.Fatalf("expected to fail create %d size file: %v", testStorageQuota*3, err)
	}
}

func TestStorageRemove(t *testing.T) {
	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	var (
		err error
	)

	if err = testStorage.Remove(); err != nil {
		t.Fatalf("failed to remove storage: %v", err)
	}
}

func TestStorageCreateWithoutQuota(t *testing.T) {
	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	var (
		err  error
		stor Storage
		trg  *os.File
	)

	if stor, err = NewStorage(Config{}); err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer stor.Remove()

	if trg, err = ioutil.TempFile(stor.Path(), ""); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := io.CopyN(trg, rand.Reader, int64(testStorageQuota*100)); err != nil {
		t.Fatalf("expected to create %d size file: %v", testStorageQuota, err)
	}
}

func TestBind(t *testing.T) {
	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	const (
		testDir = "miniverse"
	)

	var (
		err error
		src Storage
		trg string
		rst Storage
	)

	if src, err = NewStorage(Config{}); err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	if trg, err = ioutil.TempDir("", ""); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	if rst, err = Bind(trg, src.Path(), 0); err != nil {
		t.Fatalf("failed to bind: %v", err)
	}

	if err = os.Mkdir(filepath.Join(trg, testDir), os.ModePerm); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	if _, err = os.Stat(filepath.Join(src.Path(), testDir)); err != nil {
		p := filepath.Join(src.Path(), testDir)
		if os.IsNotExist(err) {
			t.Fatalf("expected that the %s exists", p)
		}
		t.Fatalf("failed to stat %s: %v", p, err)
	}

	if _, err = os.Stat(filepath.Join(trg, testDir)); err != nil {
		p := filepath.Join(trg, testDir)
		if os.IsNotExist(err) {
			t.Fatalf("expected that the %s exists", p)
		}
		t.Fatalf("failed to stat %s: %v", p, err)
	}

	if err = os.RemoveAll(trg); err == nil {
		t.Fatalf("expected that fail to remove %s", trg)
	}

	if err = src.Remove(); err != nil {
		t.Fatalf("failed to remove %s: %v", src.Path(), err)
	}

	if err = rst.Remove(); err != nil {
		t.Fatalf("failed to remove %s: %v", rst.Path(), err)
	}

	if err = os.RemoveAll(trg); err != nil {
		t.Fatalf("failed to remove %s: %v", trg, err)
	}
}
