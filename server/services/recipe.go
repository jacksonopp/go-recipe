package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
	"time"
)

type RecipeServiceErrorCode int

const (
	ErrUnknownRecipe RecipeServiceErrorCode = iota
	ErrRecipeNotFound
	ErrIngredientNotFound
	ErrIngredientConflict
	ErrInstructionNotFound
	ErrInstructionConflict
	ErrTagConflict
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
	// RECIPES
	CreateRecipe(userID uint, name, description string, ingredients []domain.IngredientDto, instructions []domain.InstructionDto) (*domain.Recipe, error)
	GetRecipeById(id uint) (*domain.Recipe, error)
	UpdateRecipe(recipeID uint, name, description string) (*domain.Recipe, error)
	DeleteRecipe(recipeID uint) error

	// INGREDIENTS
	AddIngredientToRecipe(recipeId uint, name, quantity, unit string) (*domain.Recipe, error)
	UpdateIngredient(recipeID, ingredientID uint, name, qty, unit string) (*domain.Recipe, error)
	DeleteIngredient(recipeID, ingredientID uint) error

	// INSTRUCTIONS
	AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error)
	UpdateInstruction(recipeID, instructionID uint, contents string) (*domain.Recipe, error)
	SwapInstructions(recipeID, instructionOneID, instructionTwoID uint) (*domain.Recipe, error)
	DeleteInstruction(recipeID, instructionID uint) error

	//	TAGS
	AddTagToRecipe(recipeID uint, tagId uint) (*domain.Recipe, error)
}

type recipeService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewRecipeService(db *gorm.DB) RecipeService {
	ctx := context.Background()
	return &recipeService{db: db, ctx: ctx}
}

type recipeVal struct {
	recipe *domain.Recipe
	err    error
}

// RECIPES

// CreateRecipe creates a new recipe with the given name and description.
func (r *recipeService) CreateRecipe(userID uint, name, description string, ingredients []domain.IngredientDto, instructions []domain.InstructionDto) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
		log.Println("creating recipe", name, description, userID)
		tx := r.db.Begin()
		defer recoverTx(tx)

		recipe := &domain.Recipe{
			Name:        name,
			Description: description,
			UserID:      userID,
		}

		err := tx.Create(recipe).Error
		if err != nil {
			err := fmt.Sprintf("error creating recipe: %v", err)
			log.Println(err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(ErrUnknownRecipe, err),
			}
			return
		}

		ings := make([]domain.Ingredient, len(ingredients))
		for i, ingredient := range ingredients {
			ings[i] = domain.Ingredient{
				Name:     ingredient.Name,
				Quantity: ingredient.Quantity,
				Unit:     ingredient.Unit,
				RecipeID: recipe.ID,
			}
		}
		err = tx.Create(&ings).Error
		if err != nil {
			err := fmt.Sprintf("error creating ingredients: %v", err)
			log.Println(err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(ErrUnknownRecipe, err),
			}
			return
		}

		insts := make([]domain.Instruction, len(instructions))
		for i, instruction := range instructions {
			insts[i] = domain.Instruction{
				Step:     instruction.Step,
				Contents: instruction.Contents,
				RecipeID: recipe.ID,
			}
		}
		err = tx.Create(&insts).Error
		if err != nil {
			err := fmt.Sprintf("error creating instructions: %v", err)
			log.Println(err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(ErrUnknownRecipe, err),
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			err := fmt.Sprintf("error committing transaction: %v", err)
			log.Println(err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(ErrUnknownRecipe, err),
			}
			return
		}

		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout creating recipe")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

// GetRecipeById returns the recipe with the given ID.
func (r *recipeService) GetRecipeById(id uint) (*domain.Recipe, error) {
	return getRecipeByIdWithTx(r.ctx, r.db, id)
}

// DeleteRecipe deletes the recipe with the given ID.
func (r *recipeService) DeleteRecipe(recipeID uint) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	errCh := make(chan error)

	go func() {
		defer cancel()
		tx := r.db.Begin()
		defer recoverTx(tx)

		err := tx.Delete(&domain.Recipe{}, recipeID).Error
		if err != nil {
			log.Println("error deleting recipe", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting recipe: %v", err))
			return
		}
		err = tx.Delete(&domain.Ingredient{}, "recipe_id = ?", recipeID).Error
		if err != nil {
			log.Println("error deleting ingredients", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting ingredients: %v", err))
			return
		}
		err = tx.Delete(&domain.Instruction{}, "recipe_id = ?", recipeID).Error
		if err != nil {
			log.Println("error deleting instructions", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting instructions: %v", err))
			return
		}
		err = tx.Commit().Error
		if err != nil {
			log.Println("error committing transaction", err)
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error committing transaction: %v", err))
			return
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return NewRecipeServiceError(ErrUnknownRecipe, "timeout deleting recipe")
		}
		return nil
	}
}

// INGREDIENTS

// AddIngredientToRecipe adds an ingredient to the recipe with the given ID.
func (r *recipeService) AddIngredientToRecipe(recipeID uint, name, quantity, unit string) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
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
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error creating ingredient: %v", err),
				),
			}
			return
		}
		recipe, err := getRecipeByIdWithTx(r.ctx, tx, recipeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ch <- recipeVal{
					nil,
					NewRecipeServiceError(
						ErrRecipeNotFound,
						fmt.Sprintf("recipe %d not found", recipeID),
					),
				}
				return
			}
			log.Println("error getting recipe", err)
			tx.Rollback()
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}
		if err := tx.Commit().Error; err != nil {
			log.Println("error committing transaction", err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error committing transaction: %v", err),
				),
			}
			return
		}
		ch <- recipeVal{recipe, nil}
	}()

	//return recipe, nil
	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout adding ingredient")
		} else {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
		}
	}
}

// UpdateRecipe updates the recipe with the given ID.
func (r *recipeService) UpdateRecipe(recipeID uint, name, description string) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)

	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()

		tx := r.db.Begin()

		recipe, err := getRecipeByIdWithTx(r.ctx, tx, recipeID)
		if err != nil {
			tx.Rollback()
			ch <- recipeVal{
				nil,
				err,
			}
			return
		}

		if name != "" {
			recipe.Name = name
		}
		if description != "" {
			recipe.Description = description
		}

		err = tx.Save(&recipe).Error
		if err != nil {
			tx.Rollback()
			log.Println("error saving recipe", err)
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving recipe: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error saving recipe: %v", err),
				),
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println("error committing transaction", err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error committing transaction: %v", err),
				),
			}
			return
		}
		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout updating recipe")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

// UpdateIngredient updates the ingredient with the given ID.
func (r *recipeService) UpdateIngredient(recipeID, ingredientID uint, name, qty, unit string) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
		tx := r.db.Begin()
		defer recoverTx(tx)

		var ingredient domain.Ingredient
		err := tx.First(&ingredient, ingredientID).Error
		if err != nil {
			log.Println("error getting ingredient", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				//return nil, NewRecipeServiceError(ErrIngredientNotFound, fmt.Sprintf("ingredient %d not found", ingredientID))
				ch <- recipeVal{
					nil,
					NewRecipeServiceError(
						ErrIngredientNotFound,
						fmt.Sprintf("ingredient %d not found", ingredientID),
					),
				}
				return
			}
			tx.Rollback()
			//return nil, err
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting ingredient: %v", err),
				),
			}
			return
		}

		if ingredient.RecipeID != recipeID {
			err := fmt.Sprintf("ingredient %d does not belong to recipe %d", ingredientID, recipeID)
			log.Println(err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrIngredientConflict, err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrIngredientConflict,
					err,
				),
			}
			return
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
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving ingredient: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error saving ingredient: %v", err),
				),
			}
			return
		}

		recipe, err := getRecipeByIdWithTx(r.ctx, tx, ingredient.RecipeID)
		if err != nil {
			log.Println("error getting recipe", err)
			tx.Rollback()
			//return nil, err
			ch <- recipeVal{
				nil,
				err,
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println("error getting recipe", err)
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}
		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout updating ingredient")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

// DeleteIngredient deletes the ingredient with the given ID.
func (r *recipeService) DeleteIngredient(recipeID, ingredientID uint) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	errCh := make(chan error)

	go func() {
		defer cancel()
		log.Println("deleting ingredient", recipeID, ingredientID)
		tx := r.db.Begin()
		defer recoverTx(tx)

		var ingredient domain.Ingredient

		err := tx.First(&ingredient, ingredientID).Error
		if err != nil {
			log.Println("error getting ingredient", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrIngredientNotFound, fmt.Sprintf("ingredient %d not found", ingredientID))
			return
		}
		log.Println("got ingredient", ingredient)

		if ingredient.RecipeID != recipeID {
			err := fmt.Sprintf("ingredient %d does not belong to recipe %d", ingredientID, recipeID)
			log.Println(err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrIngredientConflict, err)
			return
		}

		err = tx.Delete(&ingredient).Error
		if err != nil {
			log.Println("error deleting ingredient", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting ingredient: %v", err))
			return
		}
		tx.Commit()
		log.Println("deleted ingredient", ingredient)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return NewRecipeServiceError(ErrUnknownRecipe, "timeout deleting ingredient")
		} else {
			return nil
		}

	}
}

// INSTRUCTIONS

// AddInstructionToRecipe adds an instruction to the recipe with the given ID.
func (r *recipeService) AddInstructionToRecipe(recipeID uint, step int, contents string) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
		tx := r.db.Begin()
		defer recoverTx(tx)

		err := tx.Create(&domain.Instruction{
			Step:     step,
			Contents: contents,
			RecipeID: recipeID,
		}).Error
		if err != nil {
			log.Println("error creating instruction", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error creating instruction: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error creating instruction: %v", err),
				),
			}
			return
		}

		recipe, err := getRecipeByIdWithTx(r.ctx, tx, recipeID)
		if err != nil {
			log.Println("error getting recipe", err)
			tx.Rollback()
			//return nil, err
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println("error getting recipe", err)
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}
		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout adding instruction")
		} else {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
		}
	}
}

// UpdateInstruction updates the instruction with the given ID.
func (r *recipeService) UpdateInstruction(recipeID, instructionID uint, contents string) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
		tx := r.db.Begin()
		defer recoverTx(tx)

		var instruction domain.Instruction
		err := tx.First(&instruction, instructionID).Error
		if err != nil {
			log.Println("error getting instruction", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionID))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrInstructionNotFound,
					fmt.Sprintf("instruction %d not found", instructionID),
				),
			}
			return
		}

		if instruction.RecipeID != recipeID {
			err := fmt.Sprintf("instruction %d does not belong to recipe %d", instructionID, recipeID)
			log.Println(err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrInstructionConflict, err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrInstructionConflict,
					err,
				),
			}
			return
		}

		if contents != "" {
			instruction.Contents = contents
		}

		err = tx.Save(&instruction).Error
		if err != nil {
			log.Println("error saving instruction", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error saving instruction: %v", err),
				),
			}
			return
		}

		recipe, err := getRecipeByIdWithTx(context.Background(), tx, instruction.RecipeID)
		if err != nil {
			log.Println("error getting recipe", err)
			tx.Rollback()
			//return nil, err
			ch <- recipeVal{
				nil,
				err,
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println("error getting recipe", err)
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}

		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout updating instruction")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

// SwapInstructions swaps the positions of two instructions.
func (r *recipeService) SwapInstructions(recipeID, instructionOneID, instructionTwoID uint) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {

		tx := r.db.Begin()
		defer recoverTx(tx)

		var instructionOne domain.Instruction
		err := tx.First(&instructionOne, instructionOneID).Error
		if err != nil {
			log.Println("error getting instruction one", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionOneID))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrInstructionNotFound,
					fmt.Sprintf("instruction %d not found", instructionOneID),
				),
			}
			return
		}

		var instructionTwo domain.Instruction
		err = tx.First(&instructionTwo, instructionTwoID).Error
		if err != nil {
			log.Println("error getting instruction two", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionTwoID))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrInstructionNotFound,
					fmt.Sprintf("instruction %d not found", instructionTwoID),
				),
			}
			return
		}

		if instructionOne.RecipeID != recipeID || instructionTwo.RecipeID != recipeID {
			err := fmt.Sprintf("instructions do not belong to recipe %d", recipeID)
			log.Println(err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrInstructionConflict, err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrInstructionConflict,
					err,
				),
			}
			return
		}

		instructionOne.Step, instructionTwo.Step = instructionTwo.Step, instructionOne.Step

		err = tx.Save(&instructionOne).Error
		if err != nil {
			log.Println("error saving instruction one", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction one: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error saving instruction one: %v", err),
				),
			}
			return
		}

		err = tx.Save(&instructionTwo).Error
		if err != nil {
			log.Println("error saving instruction two", err)
			tx.Rollback()
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error saving instruction two: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error saving instruction two: %v", err),
				),
			}
			return
		}

		recipe, err := getRecipeByIdWithTx(context.Background(), tx, recipeID)
		if err != nil {
			log.Println("error getting recipe", err)
			tx.Rollback()
			//return nil, err
			ch <- recipeVal{
				nil,
				err,
			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println("error getting recipe", err)
			//return nil, NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error getting recipe: %v", err))
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error commiting transaction: %v", err),
				),
			}
		}

		ch <- recipeVal{recipe, nil}
	}()

	//return recipe, nil
	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout swapping instructions")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

// DeleteInstruction deletes the instruction with the given ID.
func (r *recipeService) DeleteInstruction(recipeID, instructionID uint) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	errCh := make(chan error)

	go func() {
		defer cancel()
		log.Println("deleting instruction", recipeID, instructionID)
		tx := r.db.Begin()
		defer recoverTx(tx)

		var instruction domain.Instruction

		err := tx.First(&instruction, instructionID).Error
		if err != nil {
			log.Println("error getting instruction", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrInstructionNotFound, fmt.Sprintf("instruction %d not found", instructionID))
			return
		}
		log.Println("got instruction", instruction)

		if instruction.RecipeID != recipeID {
			err := fmt.Sprintf("instruction %d does not belong to recipe %d", instructionID, recipeID)
			log.Println(err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrInstructionConflict, err)
			return
		}

		err = tx.Delete(&instruction).Error
		if err != nil {
			log.Println("error deleting instruction", err)
			tx.Rollback()
			errCh <- NewRecipeServiceError(ErrUnknownRecipe, fmt.Sprintf("error deleting instruction: %v", err))
			return
		}
		tx.Commit()
		log.Println("deleted instruction", instruction)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return NewRecipeServiceError(ErrUnknownRecipe, "timeout deleting instruction")
		} else {
			return nil
		}

	}
}

// TAGS

// AddTagToRecipe adds a tag to the recipe with the given ID.
// It also adds the recipe to the tag.
func (r *recipeService) AddTagToRecipe(recipeID, tagID uint) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(r.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()
		tx := r.db.Begin()
		defer recoverTx(tx)

		// Get recipe
		recipe, err := getRecipeByIdWithTx(r.ctx, tx, recipeID)
		if err != nil {
			log.Println("error getting recipe", err)
			tx.Rollback()
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}

		// Check if tag already exists in recipe
		for _, tag := range recipe.Tags {
			if tag.ID == tagID {
				err := fmt.Sprintf("tag %d already exists in recipe %d", tagID, recipeID)
				log.Println(err)
				tx.Rollback()
				ch <- recipeVal{
					nil,
					NewRecipeServiceError(
						ErrTagConflict,
						err,
					),
				}
				return
			}
		}

		// Get tag
		var tag domain.Tag
		err = tx.First(&tag, tagID).Error
		if err != nil {
			log.Println("error getting tag", err)
			tx.Rollback()
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting tag: %v", err),
				),
			}
			return
		}

		// Add tag to recipe
		err = tx.Model(&recipe).Association("Tags").Append(&tag)
		if err != nil {
			log.Println("error adding tag to recipe", err)
			tx.Rollback()
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error adding tag to recipe: %v", err),
				),
			}
			return
		}

		// Add recipe to tag
		err = tx.Model(&domain.Tag{}).
			Where("id = ?", tagID).
			Association("Recipes").
			Append(&recipe)
		if err != nil {
			log.Println("error adding recipe to tag", err)
			tx.Rollback()
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error adding recipe to tag: %v", err),
				),
			}
			return
		}

		// Commit transaction
		err = tx.Commit().Error
		if err != nil {
			log.Println("error committing transaction", err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error committing transaction: %v", err),
				),
			}
			return
		}

		ch <- recipeVal{recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout adding tag to recipe")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}

func getRecipeByIdWithTx(ctx context.Context, tx *gorm.DB, id uint) (*domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ch := make(chan recipeVal)

	go func() {
		defer cancel()

		var recipe domain.Recipe
		err := tx.
			Preload("Ingredients").
			Preload("Tags").
			First(&recipe, id).Error
		if err != nil {
			log.Println("error getting recipe", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ch <- recipeVal{
					nil,
					NewRecipeServiceError(
						ErrRecipeNotFound,
						fmt.Sprintf("recipe %d not found", id),
					),
				}
				return
			}
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting recipe: %v", err),
				),
			}
			return
		}

		var instructions []domain.Instruction
		err = tx.Order("instructions.step ASC").Find(&instructions, "recipe_id = ?", id).Error
		if err != nil {
			log.Println("error getting instructions", err)
			ch <- recipeVal{
				nil,
				NewRecipeServiceError(
					ErrUnknownRecipe,
					fmt.Sprintf("error getting instructions: %v", err),
				),
			}
			return
		}

		recipe.Instructions = instructions
		ch <- recipeVal{&recipe, nil}
	}()

	select {
	case v := <-ch:
		return v.recipe, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout getting recipe")
		}
		return nil, NewRecipeServiceError(ErrUnknownRecipe, "timeout cancelled without error")
	}
}
