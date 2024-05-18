package services

import (
	"errors"
	"fmt"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type RecipeServiceErrorCode int

const (
	ErrUnknownRecipe RecipeServiceErrorCode = iota
	ErrRecipeNotFound
	ErrIngredientNotFound
	ErrIngredientConflict
	ErrInstructionNotFound
	ErrInstructionConflict
)

type RecipeServiceError struct {
	Message string
	Code    RecipeServiceErrorCode
}

func NewRecipeServiceError(code RecipeServiceErrorCode, msg string) RecipeServiceError {
	return RecipeServiceError{
		Message: msg,
		Code:    code,
	}
}

func (e RecipeServiceError) Error() string {
	return fmt.Sprintf("recipe service error: %v", e.Message)
}

type RecipeService interface {
	CreateRecipe(userID uint, name, description string) error
	GetRecipeById(id uint) (*domain.Recipe, error)
	AddIngredientToRecipe(recipeId uint, name, quantity, unit string) (*domain.Recipe, error)
	AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error)
	UpdateRecipe(recipeID uint, name, description string) (*domain.Recipe, error)
	UpdateIngredient(recipeID, ingredientID uint, name, qty, unit string) (*domain.Recipe, error)
	DeleteIngredient(recipeID, ingredientID uint) error
	UpdateInstruction(recipeID, instructionID uint, contents string) (*domain.Recipe, error)
	SwapInstructions(recipeID, instructionOneID, instructionTwoID uint) (*domain.Recipe, error)
}

type recipeService struct {
	db *gorm.DB
}

func NewRecipeService(db *gorm.DB) RecipeService {
	return &recipeService{db: db}
}

func (r *recipeService) CreateRecipe(userID uint, name, description string) error {
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
	return getRecipeByIdWithTx(r.db, id)
}

func getRecipeByIdWithTx(tx *gorm.DB, id uint) (*domain.Recipe, error) {
	var recipe domain.Recipe
	err := tx.Preload("Ingredients").First(&recipe, id).Error
	if err != nil {
		log.Println("error getting recipe", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewRecipeServiceError(ErrRecipeNotFound, fmt.Sprintf("recipe %d not found", id))
		}
		return nil, err
	}

	var instructions []domain.Instruction
	err = tx.Order("instructions.step ASC").Find(&instructions, "recipe_id = ?", id).Error
	if err != nil {
		log.Println("error getting instructions", err)
		return nil, NewRecipeServiceError(ErrRecipeNotFound, fmt.Sprintf("recipe %d not found", id))
	}

	recipe.Instructions = instructions

	return &recipe, nil
}

func (r *recipeService) AddIngredientToRecipe(recipeID uint, name, quantity, unit string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer recoverTx(tx)

	err := tx.Create(&domain.Ingredient{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
		RecipeID: recipeID,
	}).Error
	if err != nil {
		log.Println("error creating ingredient", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error creating ingredient: %v", err))
	}
	recipe, err := getRecipeByIdWithTx(tx, recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewRecipeServiceError(ErrRecipeNotFound, fmt.Sprintf("recipe %d not found", recipeID))
		}
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
	}
	if err := tx.Commit().Error; err != nil {
		log.Println("error committing transaction", err)
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error committing transaction: %v", err))
	}

	return recipe, nil
}

func (r *recipeService) AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(&domain.Instruction{
		Step:     step,
		Contents: contents,
		RecipeID: recipeID,
	}).Error
	if err != nil {
		log.Println("error creating instruction", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error creating instruction: %v", err))
	}
	recipe, err := getRecipeByIdWithTx(tx, recipeID)
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return recipe, nil
}

func (r *recipeService) UpdateRecipe(recipeID uint, name, description string) (*domain.Recipe, error) {
	var recipe domain.Recipe

	err := r.db.First(&recipe, recipeID).Error
	if err != nil {
		log.Println("error getting recipe", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewRecipeServiceError(ErrRecipeNotFound, fmt.Sprintf("recipe %d not found", recipeID))
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
	}

	if name != "" {
		recipe.Name = name
	}
	if description != "" {
		recipe.Description = description
	}

	err = r.db.Save(&recipe).Error
	if err != nil {
		log.Println("error saving recipe", err)
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving recipe: %v", err))
	}

	return &recipe, nil
}

func (r *recipeService) UpdateIngredient(recipeID, ingredientID uint, name, qty, unit string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer recoverTx(tx)

	var ingredient domain.Ingredient
	err := tx.First(&ingredient, ingredientID).Error
	if err != nil {
		log.Println("error getting ingredient", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, NewRecipeServiceError(ErrIngredientNotFound, fmt.Sprintf("ingredient %d not found", ingredientID))
		}
		tx.Rollback()
		return nil, err
	}

	if ingredient.RecipeID != recipeID {
		err := fmt.Sprintf("ingredient %d does not belong to recipe %d", ingredientID, recipeID)
		log.Println(err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrIngredientConflict, err)
	}

	if name != "" {
		ingredient.Name = name
	}
	if qty != "" {
		ingredient.Quantity = qty
	}
	if unit != "" {
		ingredient.Unit = unit
	}

	err = tx.Save(&ingredient).Error
	if err != nil {
		log.Println("error saving ingredient", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving ingredient: %v", err))
	}

	recipe, err := getRecipeByIdWithTx(tx, ingredient.RecipeID)
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Println("error getting recipe", err)
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
	}

	return recipe, nil
}

func (r *recipeService) DeleteIngredient(recipeID, ingredientID uint) error {
	tx := r.db.Begin()
	defer recoverTx(tx)

	var ingredient domain.Ingredient

	err := tx.First(&ingredient, ingredientID).Error
	if err != nil {
		log.Println("error getting ingredient", err)
		tx.Rollback()
		return NewRecipeServiceError(ErrIngredientNotFound, fmt.Sprintf("ingredient %d not found", ingredientID))
	}
	if ingredient.RecipeID != recipeID {
		err := fmt.Sprintf("ingredient %d does not belong to recipe %d", ingredientID, recipeID)
		log.Println(err)
		tx.Rollback()
		return NewRecipeServiceError(ErrIngredientConflict, err)
	}

	err = tx.Delete(&ingredient).Error
	if err != nil {
		log.Println("error deleting ingredient", err)
		tx.Rollback()
		return NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting ingredient: %v", err))
	}

	return nil
}

func (r *recipeService) UpdateInstruction(recipeID, instructionID uint, contents string) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer recoverTx(tx)

	var instruction domain.Instruction
	err := tx.First(&instruction, instructionID).Error
	if err != nil {
		log.Println("error getting instruction", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionID))
	}

	if instruction.RecipeID != recipeID {
		err := fmt.Sprintf("instruction %d does not belong to recipe %d", instructionID, recipeID)
		log.Println(err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrInstructionConflict, err)
	}

	if contents != "" {
		instruction.Contents = contents
	}

	err = tx.Save(&instruction).Error
	if err != nil {
		log.Println("error saving instruction", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction: %v", err))
	}

	recipe, err := getRecipeByIdWithTx(tx, instruction.RecipeID)
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Println("error getting recipe", err)
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
	}

	return recipe, nil
}

func (r *recipeService) SwapInstructions(recipeID, instructionOneID, instructionTwoID uint) (*domain.Recipe, error) {
	tx := r.db.Begin()
	defer recoverTx(tx)

	var instructionOne domain.Instruction
	err := tx.First(&instructionOne, instructionOneID).Error
	if err != nil {
		log.Println("error getting instruction one", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionOneID))
	}

	var instructionTwo domain.Instruction
	err = tx.First(&instructionTwo, instructionTwoID).Error
	if err != nil {
		log.Println("error getting instruction two", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionTwoID))
	}

	if instructionOne.RecipeID != recipeID || instructionTwo.RecipeID != recipeID {
		err := fmt.Sprintf("instructions do not belong to recipe %d", recipeID)
		log.Println(err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrInstructionConflict, err)
	}

	instructionOne.Step, instructionTwo.Step = instructionTwo.Step, instructionOne.Step

	err = tx.Save(&instructionOne).Error
	if err != nil {
		log.Println("error saving instruction one", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction one: %v", err))
	}

	err = tx.Save(&instructionTwo).Error
	if err != nil {
		log.Println("error saving instruction two", err)
		tx.Rollback()
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction two: %v", err))
	}

	recipe, err := getRecipeByIdWithTx(tx, recipeID)
	if err != nil {
		log.Println("error getting recipe", err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Println("error getting recipe", err)
		return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
	}

	return recipe, nil
}

func recoverTx(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}
