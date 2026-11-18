package jsonutils

import (
	"github.com/username/go-utils/pkg/jsonutils/models"
	"strings"
	"testing"
	"time"
)

func TestToDomain(t *testing.T) {
	t.Run("FullData", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		age := 30
		street := "123 Main St"
		city := "New York"
		theme := "dark"
		notif := true
		createdAt := "2023-01-01T00:00:00Z"
		parsedTime, _ := time.Parse(time.RFC3339, createdAt)

		dto := &models.UserDTO{
			Name:  &name,
			Email: &email,
			Age:   &age,
			Address: &models.AddressDTO{
				Street: &street,
				City:   &city,
			},
			Settings: &models.SettingsDTO{
				NotificationsEnabled: &notif,
				Theme:                &theme,
			},
			CreatedAt: &createdAt,
		}

		domain, err := ToDomain(dto)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if domain.Name != name {
			t.Errorf("expected Name %s, got %s", name, domain.Name)
		}
		if domain.Email != email {
			t.Errorf("expected Email %s, got %s", email, domain.Email)
		}
		if domain.Age != age {
			t.Errorf("expected Age %d, got %d", age, domain.Age)
		}
		if domain.Address.Street != street {
			t.Errorf("expected Street %s, got %s", street, domain.Address.Street)
		}
		if domain.Address.City != city {
			t.Errorf("expected City %s, got %s", city, domain.Address.City)
		}
		if domain.Settings.NotificationsEnabled != notif {
			t.Errorf("expected NotificationsEnabled %v, got %v", notif, domain.Settings.NotificationsEnabled)
		}
		if domain.Settings.Theme != theme {
			t.Errorf("expected Theme %s, got %s", theme, domain.Settings.Theme)
		}
		if !domain.CreatedAt.Equal(parsedTime) {
			t.Errorf("expected CreatedAt %v, got %v", parsedTime, domain.CreatedAt)
		}
	})

	t.Run("NullFields", func(t *testing.T) {
		dto := &models.UserDTO{}

		domain, err := ToDomain(dto)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if domain.Name != "" {
			t.Errorf("expected empty Name, got %s", domain.Name)
		}
		if domain.Email != "" {
			t.Errorf("expected empty Email, got %s", domain.Email)
		}
		if domain.Age != 0 {
			t.Errorf("expected Age 0, got %d", domain.Age)
		}
		if domain.Address.Street != "" {
			t.Errorf("expected empty Street, got %s", domain.Address.Street)
		}
		if domain.Settings.Theme != "light" {
			t.Errorf("expected default Theme 'light', got %s", domain.Settings.Theme)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		tests := []struct {
			name          string
			dto           *models.UserDTO
			expectedError string
		}{
			{
				name: "NegativeAge",
				dto: &models.UserDTO{
					Age: ptrInt(-1),
				},
				expectedError: "age must be between 0 and 150",
			},
			{
				name: "TooOld",
				dto: &models.UserDTO{
					Age: ptrInt(151),
				},
				expectedError: "age must be between 0 and 150",
			},
			{
				name: "InvalidEmail",
				dto: &models.UserDTO{
					Email: ptrString("invalid-email"),
				},
				expectedError: "invalid email format",
			},
			{
				name: "TooLongName",
				dto: &models.UserDTO{
					Name: ptrString(strings.Repeat("a", 101)),
				},
				expectedError: "name exceeds 100 characters",
			},
			{
				name: "InvalidTheme",
				dto: &models.UserDTO{
					Settings: &models.SettingsDTO{
						Theme: ptrString("blue"),
					},
				},
				expectedError: "theme must be light or dark",
			},
			{
				name: "InvalidCreatedAt",
				dto: &models.UserDTO{
					CreatedAt: ptrString("invalid-date"),
				},
				expectedError: "invalid created_at format",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := ToDomain(tt.dto)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			})
		}
	})
}

func ptrInt(i int) *int {
	return &i
}

func ptrString(s string) *string {
	return &s
}
