package models

import "gorm.io/gorm"

type Text struct {
	gorm.Model
	To        string `gorm:"required"`
	MessageID uint   `gorm:"foreignkey:Message"`
	Content   string
	Status    bool `gorm:"default:true"`
}
