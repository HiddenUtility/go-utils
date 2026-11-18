package jsonutils

import (
	"github.com/username/go-utils/pkg/jsonutils/models"
	"net/mail"
	"time"
)

// ToDomain converts a UserDTO to a User domain model, applying defaults and validation.
func ToDomain(dto *models.UserDTO) (*models.User, error) {
	if dto == nil {
		return nil, nil
	}

	domain := &models.User{
		Name:      "",
		Email:     "",
		Age:       0,
		Address:   models.Address{Street: "", City: ""},
		Settings:  models.Settings{NotificationsEnabled: true, Theme: "light"},
		CreatedAt: time.Time{},
	}

	// Apply values and handle defaults
	if dto.Name != nil {
		domain.Name = *dto.Name
	}
	if dto.Email != nil {
		domain.Email = *dto.Email
	}
	if dto.Age != nil {
		domain.Age = *dto.Age
	}
	if dto.Address != nil {
		if dto.Address.Street != nil {
			domain.Address.Street = *dto.Address.Street
		}
		if dto.Address.City != nil {
			domain.Address.City = *dto.Address.City
		}
	}
	if dto.Settings != nil {
		if dto.Settings.NotificationsEnabled != nil {
			domain.Settings.NotificationsEnabled = *dto.Settings.NotificationsEnabled
		}
		if dto.Settings.Theme != nil {
			domain.Settings.Theme = *dto.Settings.Theme
		}
	}
	if dto.CreatedAt != nil {
		t, err := time.Parse(time.RFC3339, *dto.CreatedAt)
		if err != nil {
			return nil, ValidationError{Message: "invalid created_at format"}
		}
		domain.CreatedAt = t
	}

	// Validation
	if len(domain.Name) > 100 {
		return nil, ValidationError{Message: "name exceeds 100 characters"}
	}
	if domain.Email != "" {
		if _, err := mail.ParseAddress(domain.Email); err != nil {
			return nil, ValidationError{Message: "invalid email format"}
		}
	}
	if domain.Age < 0 || domain.Age > 150 {
		return nil, ValidationError{Message: "age must be between 0 and 150"}
	}
	if domain.Settings.Theme != "" && domain.Settings.Theme != "light" && domain.Settings.Theme != "dark" {
		return nil, ValidationError{Message: "theme must be light or dark"}
	}

	return domain, nil
}
