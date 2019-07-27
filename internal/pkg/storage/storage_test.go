package storage_test

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
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
