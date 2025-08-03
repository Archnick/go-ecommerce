package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TotalAmount     float64     `gorm:"type:decimal(10,2);not null"`
	Status          string      `gorm:"type:varchar(50);deafault:'pending';not null"`
	ShippingAddress string      `gorm:"type:text;not null"`
	UserID          uint        `gorm:"type:not null"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID"`
}
