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
	// UserID is the ID of the user who created the recipe.
	UserID       uint          `json:"user_id"`
	Ingredients  []Ingredient  `gorm:"foreignKey:RecipeID"`
	Instructions []Instruction `gorm:"foreignKey:RecipeID"`
	Tags         []*Tag        `gorm:"many2many:recipe_tags"`
}

// RecipeDto is a DTO for a Recipe.
type RecipeDto struct {
	ID           uint        `json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Ingredients  []Dto       `json:"ingredients"`
	Instructions []Dto       `json:"instructions"`
	Tags         []simpleTag `json:"tags"`
}

type simpleTag struct {
	ID  uint   `json:"id"`
	Tag string `json:"tag"`
}

// ToDto converts a Recipe to a RecipeDto.
func (r *Recipe) ToDto() Dto {
	ingredients := make([]Dto, len(r.Ingredients))
	for i, ingredient := range r.Ingredients {
		ingredients[i] = ingredient.ToDto()
	}

	instructions := make([]Dto, len(r.Instructions))
	for i, instruction := range r.Instructions {
		instructions[i] = instruction.ToDto()
	}

	tags := make([]simpleTag, len(r.Tags))
	for i, tag := range r.Tags {
		tags[i] = simpleTag{
			ID:  tag.ID,
			Tag: tag.Tag,
		}
	}

	return RecipeDto{
		ID:           r.ID,
		CreatedAt:    r.CreatedAt,
		Name:         r.Name,
		Description:  r.Description,
		Ingredients:  ingredients,
		Instructions: instructions,
		Tags:         tags,
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

type IngredientDto struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

// ToDto converts an Ingredient to a Dto.
func (i *Ingredient) ToDto() Dto {
	return IngredientDto{
		ID:       i.ID,
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

type InstructionDto struct {
	ID       uint   `json:"id"`
	Step     int    `json:"step"`
	Contents string `json:"contents"`
}

// ToDto converts an Instruction to a Dto.
func (i *Instruction) ToDto() Dto {
	return InstructionDto{
		ID:       i.ID,
		Step:     i.Step,
		Contents: i.Contents,
	}
}
