package io

import (
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {
	src := "/tmp/test_src.txt"
	dst := "/tmp/test_dst.txt"

	// Ensure src exists
	err := os.WriteFile(src, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("failed to create src: %v", err)
	}

	err = CopyFile(src, dst)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read dst: %v", err)
	}

	if string(content) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(content))
	}

	// Cleanup
	os.Remove(src)
	os.Remove(dst)
}
