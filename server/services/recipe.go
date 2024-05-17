package services

import (
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type RecipeService interface {
	CreateRecipe(name, description string, userID uint) error
	GetRecipeById(id uint) (*domain.Recipe, error)
	AddIngredientToRecipe(recipeId uint, name, quantity, unit string) (*domain.Recipe, error)
	AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error)
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

func (r *recipeService) GetRecipeById(id uint) (*domain.Recipe, error) {
	var recipe domain.Recipe
	err := r.db.Preload("Ingredients").Preload("Instructions").First(&recipe, id).Error
	if err != nil {
		log.Println("error getting recipe", err)
		return nil, err
	}

	return &recipe, nil
}

func (r *recipeService) AddIngredientToRecipe(recipeID uint, name, quantity, unit string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var recipe domain.Recipe

	err := tx.Create(&domain.Ingredient{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
		RecipeID: recipeID,
	}).Error
	if err != nil {
		log.Println("error creating ingredient", err)
		tx.Rollback()
		return nil, err
	}
	err = tx.Preload("Ingredients").Preload("Instructions").First(&recipe, recipeID).Error
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		log.Println("error committing transaction", err)
		return nil, err
	}

	return &recipe, nil
}

func (r *recipeService) AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var recipe domain.Recipe

	err := tx.Create(&domain.Instruction{
		Step:     step,
		Contents: contents,
		RecipeID: recipeID,
	}).Error
	if err != nil {
		log.Println("error creating instruction", err)
		tx.Rollback()
		return nil, err
	}
	err = tx.Preload("Ingredients").Preload("Instructions").First(&recipe, recipeID).Error
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return &recipe, nil
}
