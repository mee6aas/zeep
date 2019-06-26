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

	testStorageQuota = uint32(1024 * 4) // 4KiB
	testStoragePath  string
)

func TestStorageCreate(t *testing.T) {
	var (
		err error
		res CreateCreatedBody
	)

	if res, err = Create(CreateConfig{
		Size: testStorageQuota * 2,
	}); err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	testStoragePath = res.Path
}

func TestStorageQuota(t *testing.T) {
	var (
		err error
		trg *os.File
	)

	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	if trg, err = ioutil.TempFile(testStoragePath, ""); err != nil {
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
	var (
		err error
	)

	if testStorageCreateFailed {
		t.Skipf("TestStorageCreate failed")
	}

	if err = Remove(testStoragePath); err != nil {
		t.Fatalf("failed to remove storage: %v", err)
	}
}
