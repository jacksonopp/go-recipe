package domain

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Tag     string    `gorm:"unique;not null;index"`
	Recipes []*Recipe `gorm:"many2many:recipe_tags"`
}

type TagDto struct {
	Tag     string      `json:"tag"`
	Recipes []RecipeDto `json:"recipes"`
}

func (t *Tag) ToDto() Dto {
	recipes := make([]RecipeDto, len(t.Recipes))
	for i, recipe := range t.Recipes {
		recipes[i] = recipe.ToDto().(RecipeDto)
	}

	return TagDto{
		Tag:     t.Tag,
		Recipes: recipes,
	}
}
