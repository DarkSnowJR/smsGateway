package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Name     string
	Email    string  `gorm:"unique"`
	Balance  float64 `gorm:"default:1000"`
	SendingRate float64 `json:"sending_rate" gorm:"column:sending_rate"`
}
