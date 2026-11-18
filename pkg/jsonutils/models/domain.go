package models

import "time"

// User is the domain model for a User.
type User struct {
	Name      string
	Email     string
	Age       int
	Address   Address
	Settings  Settings
	CreatedAt time.Time
}

// Address is the domain model for an Address.
type Address struct {
	Street string
	City   string
}

// Settings is the domain model for Settings.
type Settings struct {
	NotificationsEnabled bool
	Theme                string
}
