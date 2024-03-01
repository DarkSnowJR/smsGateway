package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Name     string
	Email    string  `gorm:"unique"`
	Balance  float64 `gorm:"default:1000"`
}
