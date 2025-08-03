package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user record in the database.
// We embed gorm.Model to get the ID, CreatedAt, etc. fields for free.
type User struct {
	gorm.Model
	FirstName             string `gorm:"type:varchar(100)"`
	LastName              string `gorm:"type:varchar(100)"`
	Email                 string `gorm:"type:varchar(255);unique;not null"`
	Password              string
	ProfileImageURL       string `gorm:"type:varchar(255)"`
	Role                  string `gorm:"type:varchar(20);default:'customer'"`
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	Shop                  Shop     `gorm:"foreignKey:UserID"`
	Orders                []Order  `gorm:"foreignKey:UserID"`
	Reviews               []Review `gorm:"foreignKey:UserID"`
}
