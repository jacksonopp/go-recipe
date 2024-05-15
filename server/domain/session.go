package domain

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	UserID    uint
	Token     string    `gorm:"index;unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
