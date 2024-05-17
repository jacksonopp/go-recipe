package domain

import (
	"gorm.io/gorm"
	"time"
)

// Session represents an active or expired session.
// It is associated with a single User
type Session struct {
	gorm.Model
	// UserID is the ID of the user who created the session.
	UserID    uint
	Token     string    `gorm:"index;unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
