package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	src := "/tmp/test_dir_src"
	dst := "/tmp/test_dir_dst"

	// Ensure src exists (done in setup)
	
	err := CopyDir(src, dst)
	if err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Check if dst exists and has the same structure
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		dstInfo, err := os.Stat(dstPath)
		if err != nil {
			t.Errorf("dstPath %s does not exist", dstPath)
			return err
		}

		if dstInfo.IsDir() != info.IsDir() {
			t.Errorf("type mismatch for %s: expected dir=%v, got dir=%v", path, info.IsDir(), dstInfo.IsDir())
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read src file %s: %v", path, err)
			}
			dstContent, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("failed to read dst file %s: %v", dstPath, err)
			}
			if string(content) != string(dstContent) {
				t.Errorf("content mismatch for %s", path)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Cleanup
	os.RemoveAll(src)
	os.RemoveAll(dst)
}
