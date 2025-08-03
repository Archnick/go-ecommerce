package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ImageURL  string `gorm:"type:varchar(255);not null"`
	AltText   string `gorm:"type:varchar(255)"`
	ProductID uint   `gorm:"not null"`
}
