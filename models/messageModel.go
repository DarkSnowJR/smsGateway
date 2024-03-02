package models

import (
	"github.com/lib/pq" // Import the pq package for PostgreSQL types
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	To      pq.StringArray `gorm:"type:text[]"`
	UserID  uint           `gorm:"foreignkey:User"`
	Content string
	Status  bool `gorm:"default:false"`
}
