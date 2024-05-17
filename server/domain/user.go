package domain

import "gorm.io/gorm"

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Salt     string    `gorm:"not null"`
	Sessions []Session `gorm:"foreignKey:UserID"`
	Recipes  []Recipe  `gorm:"foreignKey:UserID"`
}

// UserDto is a DTO for a User.
type UserDto struct {
	Username string   `json:"username"`
	Recipe   []Recipe `json:"recipes"`
}

// ToDto converts a User to a UserDto.
func (u *User) ToDto() Dto {
	return UserDto{
		Username: u.Username,
		Recipe:   u.Recipes,
	}
}
