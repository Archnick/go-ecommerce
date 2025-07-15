package models

import "gorm.io/gorm"

// User represents a user record in the database.
// We embed gorm.Model to get the ID, CreatedAt, etc. fields for free.
type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}
