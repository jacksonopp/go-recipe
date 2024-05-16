package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password"`
	Salt     string    `gorm:"not null" json:"salt"`
	Sessions []Session `gorm:"foreignKey:UserID" json:"sessions"`
}

type UserDto struct {
	Username string `json:"username"`
}

func (u *User) ToDto() Dto {
	return UserDto{
		Username: u.Username,
	}
}
