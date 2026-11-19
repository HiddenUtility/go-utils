package jsonutils

import (
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("JsonError", func(t *testing.T) {
		err := JsonError{Message: "test json error"}
		if err.Error() != "test json error" {
			t.Errorf("expected 'test json error', got '%s'", err.Error())
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := ValidationError{Message: "test validation error"}
		if err.Error() != "test validation error" {
			t.Errorf("expected 'test validation error', got '%s'", err.Error())
		}
	})
}
