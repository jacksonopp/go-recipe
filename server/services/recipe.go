package services

import (
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type RecipeService interface {
	CreateRecipe(name, description string, userID uint) error
}

type recipeService struct {
	db *gorm.DB
}

func NewRecipeService(db *gorm.DB) RecipeService {
	return &recipeService{db: db}
}

func (r *recipeService) CreateRecipe(name, description string, userID uint) error {
	log.Println("creating recipe", name, description, userID)
	err := r.db.Create(&domain.Recipe{
		Name:        name,
		Description: description,
		UserID:      userID,
	}).Error
	if err != nil {
		log.Println("error creating recipe", err)
		return err
	}
	return nil
}
