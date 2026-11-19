package models

// UserDTO is the Data Transfer Object for a User.
type UserDTO struct {
	Name      *string      `json:"name"`
	Email     *string      `json:"email"`
	Age       *int         `json:"age"`
	Address   *AddressDTO  `json:"address"`
	Settings  *SettingsDTO `json:"settings"`
	CreatedAt *string      `json:"created_at"`
}

// AddressDTO is the Data Transfer Object for an Address.
type AddressDTO struct {
	Street *string `json:"street"`
	City   *string `json:"city"`
}

// SettingsDTO is the Data Transfer Object for Settings.
type SettingsDTO struct {
	NotificationsEnabled *bool   `json:"notifications_enabled"`
	Theme                *string `json:"theme"`
}
