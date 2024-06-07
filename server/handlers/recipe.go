package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type RecipeHandler struct {
	r             fiber.Router
	db            *gorm.DB
	recipeService services.RecipeService
}

func NewRecipeHandler(r fiber.Router, db *gorm.DB) *RecipeHandler {
	subpath := r.Group("/recipe")
	recipeService := services.NewRecipeService(db)

	return &RecipeHandler{r: subpath, db: db, recipeService: recipeService}
}

func (h *RecipeHandler) RegisterRoutes() {
	// RECIPES
	h.r.Post("/", AuthMiddleware(h.db), h.createRecipe)
	h.r.Get("/:id", h.getRecipeById)
	h.r.Patch("/:id", AuthMiddleware(h.db), h.updateRecipe)
	h.r.Delete("/:id", AuthMiddleware(h.db), h.deleteRecipe)

	// INGREDIENTS
	h.r.Post("/:id/ingredient", AuthMiddleware(h.db), h.createIngredient)
	h.r.Patch("/:id/ingredient/:ingredientId", AuthMiddleware(h.db), h.updateIngredient)
	h.r.Delete("/:id/ingredient/:ingredientId", AuthMiddleware(h.db), h.deleteIngredient)
	// TODO swap ingredient

	// INSTRUCTIONS
	h.r.Post("/:id/instruction", AuthMiddleware(h.db), h.createInstruction)
	h.r.Patch("/:id/instruction/:instructionId", AuthMiddleware(h.db), h.updateInstruction)
	h.r.Patch("/:id/instruction/:instructionOneId/:instructionTwoId", AuthMiddleware(h.db), h.swapInstructions)
	h.r.Delete("/:id/instruction/:instructionId", AuthMiddleware(h.db), h.deleteInstruction)
	//	TODO delete instruction

	//	TAGS
	h.r.Patch("/:recipeId/tag/:tagId", AuthMiddleware(h.db), h.addTagToRecipe)
	h.r.Delete("/:recipeId/tag/:tagId", AuthMiddleware(h.db), h.removeTagFromRecipe)
}

// RECIPES

// POST /recipe/
func (h *RecipeHandler) createRecipe(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	recipe := struct {
		Name         string                  `json:"name"`
		Description  string                  `json:"description"`
		Ingredients  []domain.IngredientDto  `json:"ingredients"`
		Instructions []domain.InstructionDto `json:"instructions"`
		UserID       uint                    `json:"user_id"` // <- this is populated by the middleware
	}{}

	if err := c.BodyParser(&recipe); err != nil {
		log.Printf("body content: %v", string(c.Body()))
		log.Printf("error parsing body: %v", err)
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if recipe.Name == "" {
		err := UnprocessableEntity(map[string]string{"name": "name is required"})
		return SendError(c, err)
	}

	recipe.UserID = user.ID

	r, err := h.recipeService.CreateRecipe(recipe.UserID, recipe.Name, recipe.Description, recipe.Ingredients, recipe.Instructions)
	if err != nil {
		log.Println("error creating recipe", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(r.ToDto())
}

// GET /recipe/:id
func (h *RecipeHandler) getRecipeById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return SendError(c, BadRequest("id is required"))
	}

	recipeId, err := strconv.Atoi(id)
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	recipe, err := h.recipeService.GetRecipeById(uint(recipeId))
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())

	//recipe, err := h.recipeService.GetRecipeById(id)
	//if err
}

// PATCH /recipe/:id
func (h *RecipeHandler) updateRecipe(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	r := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	if err := c.BodyParser(&r); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if r.Name == "" && r.Description == "" {
		err := UnprocessableEntity(map[string]string{"error": "name or description is required"})
		return SendError(c, err)
	}

	recipe, err := h.recipeService.UpdateRecipe(user.ID, uint(id), r.Name, r.Description)
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		if errors.Is(err, services.ErrUnauthorized) {
			return SendError(c, Unauthorized())
		}
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:id
func (h *RecipeHandler) deleteRecipe(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	err = h.recipeService.DeleteRecipe(user.ID, uint(id))
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			return SendError(c, Unauthorized())
		}
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// INGREDIENTS

// POST /recipe/:id/ingredient
func (h *RecipeHandler) createIngredient(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	ingredient := struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
		Unit     string `json:"unit"`
	}{}

	if err := c.BodyParser(&ingredient); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if ingredient.Name == "" {
		err := UnprocessableEntity(map[string]string{"name": "name is required"})
		return SendError(c, err)
	}

	recipe, err := h.recipeService.AddIngredientToRecipe(user.ID, uint(id), ingredient.Name, ingredient.Quantity, ingredient.Unit)
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		log.Println("error creating ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/ingredient/:ingredientId
func (h *RecipeHandler) updateIngredient(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}
	recipeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	ingredientID, err := strconv.Atoi(c.Params("ingredientId"))
	if err != nil {
		return SendError(c, BadRequest("ingredientId must be an integer"))
	}

	ingredient := struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
		Unit     string `json:"unit"`
	}{}

	if err := c.BodyParser(&ingredient); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if ingredient.Name == "" && ingredient.Quantity == "" && ingredient.Unit == "" {
		err := UnprocessableEntity(map[string]string{"error": "name, quantity, or unit is required"})
		return SendError(c, err)
	}

	recipe, err := h.recipeService.UpdateIngredient(user.ID, uint(recipeID), uint(ingredientID), ingredient.Name, ingredient.Quantity, ingredient.Unit)
	if err != nil {
		log.Println("error updating ingredient", err)
		if errors.Is(err, services.ErrUnauthorized) {
			return SendError(c, Unauthorized())
		}
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:id/ingredient/:ingredientId
func (h *RecipeHandler) deleteIngredient(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	log.Println("delete ingredient")
	recipeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	ingredientID, err := strconv.Atoi(c.Params("ingredientId"))
	if err != nil {
		return SendError(c, BadRequest("ingredientId must be an integer"))
	}

	log.Println(recipeID, ingredientID)

	err = h.recipeService.DeleteIngredient(user.ID, uint(recipeID), uint(ingredientID))
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		if errors.Is(err, services.ErrIngredientConflict) {
			return SendError(c, Conflict(map[string]string{"error": "ingredient does not belong to recipe"}))
		}
		if errors.Is(err, services.ErrUnauthorized) {
			return SendError(c, Unauthorized())
		}
		log.Println("error deleting ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// INSTRUCTIONS

// POST /recipe
func (h *RecipeHandler) createInstruction(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	instruction := struct {
		Step     int    `json:"step"`
		Contents string `json:"contents"`
	}{}

	if err := c.BodyParser(&instruction); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if instruction.Contents == "" {
		err := UnprocessableEntity(map[string]string{"contents": "contents is required"})
		return SendError(c, err)
	}

	recipe, err := h.recipeService.AddInstructionToRecipe(user.ID, uint(id), instruction.Step, instruction.Contents)
	if err != nil {
		log.Println("error creating instruction", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/instruction/:instructionId
func (h *RecipeHandler) updateInstruction(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}
	recipeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	instructionID, err := strconv.Atoi(c.Params("instructionId"))
	if err != nil {
		return SendError(c, BadRequest("instructionId must be an integer"))
	}

	instruction := struct {
		Contents string `json:"contents"`
	}{}

	if err := c.BodyParser(&instruction); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if instruction.Contents == "" {
		err := UnprocessableEntity(map[string]string{"error": "contents is required"})
		return SendError(c, err)
	}

	recipe, err := h.recipeService.UpdateInstruction(user.ID, uint(recipeID), uint(instructionID), instruction.Contents)
	if err != nil {
		log.Println("error updating instruction", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/instruction/:instructionOneId/:instructionTwoId
func (h *RecipeHandler) swapInstructions(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}
	recipeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	instructionOneID, err := strconv.Atoi(c.Params("instructionOneId"))
	if err != nil {
		return SendError(c, BadRequest("instructionOneId must be an integer"))
	}

	instructionTwoID, err := strconv.Atoi(c.Params("instructionTwoId"))
	if err != nil {
		return SendError(c, BadRequest("instructionTwoId must be an integer"))
	}

	recipe, err := h.recipeService.SwapInstructions(user.ID, uint(recipeID), uint(instructionOneID), uint(instructionTwoID))
	if err != nil {
		log.Println("error swapping instructions", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:id/instruction/:instructionId
func (h *RecipeHandler) deleteInstruction(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	instructionId, err := strconv.Atoi(c.Params("instructionId"))
	if err != nil {
		return SendError(c, BadRequest("instructionId must be an integer"))
	}

	err = h.recipeService.DeleteInstruction(user.ID, uint(id), uint(instructionId))
	if err != nil {
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// TAGS
// PATCH /recipe/:recipeId/tag/:tagId
func (h *RecipeHandler) addTagToRecipe(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}

	recipeId, err := strconv.Atoi(c.Params("recipeId"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	tagId, err := strconv.Atoi(c.Params("tagId"))
	if err != nil {
		return SendError(c, BadRequest("invalid request body"))
	}

	recipe, err := h.recipeService.AddTagToRecipe(user.ID, uint(recipeId), uint(tagId))
	if err != nil {
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:recipeId/tag/:tagId
func (h *RecipeHandler) removeTagFromRecipe(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return SendError(c, err.(APIError))
	}
	recipeId, err := strconv.Atoi(c.Params("recipeId"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	tagId, err := strconv.Atoi(c.Params("tagId"))
	if err != nil {
		return SendError(c, BadRequest("invalid request body"))
	}

	err = h.recipeService.RemoveTagFromRecipe(user.ID, uint(recipeId), uint(tagId))
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		if errors.Is(err, services.ErrTagNotFound) {
			return SendError(c, NotFound(map[string]string{"error": "tag not found"}))
		}
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func getUserFromLocals(c *fiber.Ctx) (domain.User, error) {
	var user domain.User
	if u, ok := c.Locals("user").(*domain.User); !ok {
		log.Println("no user")
		return domain.User{}, Unauthorized()
	} else {
		user = *u
	}
	return user, nil
}
