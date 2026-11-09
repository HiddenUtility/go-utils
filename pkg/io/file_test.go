package io

import (
	"os"
	"testing"
)

func TestWriteStringToFile(t *testing.T) {
	tmpFile := "test_file.txt"
	content := "Hello, Go Utils!"

	// Cleanup after test
	defer os.Remove(tmpFile)

	err := WriteStringToFile(tmpFile, content)
	if err != nil {
		t.Fatalf("WriteStringToFile failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected %q, got %q", content, string(data))
	}
}
