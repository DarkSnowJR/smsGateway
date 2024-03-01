package models

import (
	"github.com/lib/pq" // Import the pq package for PostgreSQL types
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	To      pq.StringArray `gorm:"type:varchar(255)[]"`
	UserID  uint           `gorm:"foreignkey:User"`
	Content string
	Status  bool `gorm:"default:false"`
}
