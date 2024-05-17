package domain

import (
	"gorm.io/gorm"
	"time"
)

// Recipe represents a recipe in the system.
type Recipe struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Ingredients []Ingredient `gorm:"foreignKey:RecipeID"`
	// UserID is the ID of the user who created the recipe.
	UserID       uint          `json:"user_id"`
	Instructions []Instruction `gorm:"foreignKey:RecipeID"`
}

// RecipeDto is a DTO for a Recipe.
type RecipeDto struct {
	ID           uint          `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
}

// ToDto converts a Recipe to a RecipeDto.
func (r *Recipe) ToDto() Dto {
	return RecipeDto{
		ID:          r.ID,
		CreatedAt:   r.CreatedAt,
		Name:        r.Name,
		Description: r.Description,
		Ingredients: r.Ingredients,
	}
}

// Ingredient represents an ingredient in a Recipe.
// A Recipe can have many Ingredients.
type Ingredient struct {
	gorm.Model
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
	// RecipeID is the ID of the recipe that this Ingredient belongs to.
	RecipeID uint `json:"recipe_id"`
}

// ToDto converts an Ingredient to a Dto.
func (i *Ingredient) ToDto() Dto {
	return Ingredient{
		Name:     i.Name,
		Quantity: i.Quantity,
		Unit:     i.Unit,
	}
}

// Instruction represents an instruction in a Recipe.
// A Recipe can have many Instructions.
type Instruction struct {
	gorm.Model
	Step     int    `json:"step"`
	Contents string `json:"contents"`
	// RecipeID is the ID of the recipe that this Instruction belongs to.
	RecipeID uint `json:"recipe_id"`
}

// ToDto converts an Instruction to a Dto.
func (i *Instruction) ToDto() Dto {
	return Instruction{
		Step:     i.Step,
		Contents: i.Contents,
	}
}
