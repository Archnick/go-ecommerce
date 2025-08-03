package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	Quantity  int     `gorm:"type:integer;not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
}
