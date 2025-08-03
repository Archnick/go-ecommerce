package models

import (
	"gorm.io/gorm"
)

type Shop struct {
	gorm.Model
	Name         string    `gorm:"type:varchar(150);unique;not null"`
	Description  string    `gorm:"type:text"`
	ShopImageUrl string    `gorm:"type:varchar(255)"`
	UserID       uint      `gorm:"type:int;not null"`
	Products     []Product `gorm:"foreignKey:ShopID"`
}
