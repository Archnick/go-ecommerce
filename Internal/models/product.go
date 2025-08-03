package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"type:decimal(10,2);not null"`
	Stock       int            `gorm:"type:integer;default:0"`
	ShopID      uint           `gorm:"not null"`
	CategoryID  uint           `gorm:"not null"`
	Images      []ProductImage `gorm:"foreignKey:ProductID"`
	Reviews     []Review       `gorm:"foreignKey:ProductID"`
}
