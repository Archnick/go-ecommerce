package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Rating    int    `gorm:"type:integer;not null"`
	Comment   string `gorm:"type:text"`
	UserID    uint   `gorm:"not null"`
	ProductID uint   `gorm:"not null"`
}
