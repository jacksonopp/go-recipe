package domain

import (
	"gorm.io/gorm"
	"time"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null;uniqueIndex:idx_username"`
	Password string    `gorm:"not null"`
	Salt     string    `gorm:"not null"`
	Sessions []Session `gorm:"foreignKey:UserID"`
	Recipes  []Recipe  `gorm:"foreignKey:UserID"`
}

// UserDto is a DTO for a User.
type UserDto struct {
	ID       uint            `json:"id"`
	Username string          `json:"username"`
	Recipe   []userRecipeDto `json:"recipes"`
}

// ToDto converts a User to a UserDto.
func (u *User) ToDto() Dto {

	recipes := make([]userRecipeDto, len(u.Recipes))

	for i, recipe := range u.Recipes {
		recipes[i] = userRecipeDto{
			ID:        recipe.ID,
			CreatedAt: recipe.CreatedAt,
			Name:      recipe.Name,
		}
	}

	return UserDto{
		ID:       u.ID,
		Username: u.Username,
		Recipe:   recipes,
	}
}

// userRecipeDto is a DTO for a Recipe in a UserDto.
// It does not include the recipe's ingredients or instructions.
type userRecipeDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}
