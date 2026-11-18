package jsonutils

import (
	"github.com/username/go-utils/pkg/jsonutils/models"
	"path/filepath"
	"testing"
)

func TestJSONFunctions(t *testing.T) {
	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: models.Address{
			Street: "123 Main St",
			City:   "New York",
		},
		Settings: models.Settings{
			NotificationsEnabled: true,
			Theme:                "dark",
		},
	}

	t.Run("ToJSON", func(t *testing.T) {
		jsonStr, err := ToJSON(user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if jsonStr == "" {
			t.Fatal("expected non-empty JSON string")
		}
	})

	t.Run("FromJSON", func(t *testing.T) {
		jsonStr, _ := ToJSON(user)
		decodedUser, err := FromJSON[models.User](jsonStr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if decodedUser.Name != user.Name {
			t.Errorf("expected Name %s, got %s", user.Name, decodedUser.Name)
		}
	})

	t.Run("ToJSONPretty", func(t *testing.T) {
		jsonStr, err := ToJSONPretty(user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(jsonStr) == 0 {
			t.Fatal("expected non-empty JSON string")
		}
	})

	t.Run("FileIO", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "user.json")

		// WriteFile
		err := WriteFile(filePath, user)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		// ReadFile
		readUser, err := ReadFile[models.User](filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if readUser.Name != user.Name {
			t.Errorf("expected Name %s, got %s", user.Name, readUser.Name)
		}

		// WriteFilePretty
		prettyPath := filepath.Join(tmpDir, "user_pretty.json")
		err = WriteFilePretty(prettyPath, user)
		if err != nil {
			t.Fatalf("failed to write pretty file: %v", err)
		}

		// Read back pretty file
		readPrettyUser, err := ReadFile[models.User](prettyPath)
		if err != nil {
			t.Fatalf("failed to read pretty file: %v", err)
		}
		if readPrettyUser.Name != user.Name {
			t.Errorf("expected Name %s, got %s", user.Name, readPrettyUser.Name)
		}
	})
}
