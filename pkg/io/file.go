// Package io provides utility functions for file operations.
package io

import (
	"fmt"
	"os"
)

// WriteStringToFile writes a string to a file with specified permissions.
func WriteStringToFile(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write string to file: %w", err)
	}
	return nil
}
