package handlers

import (
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
	// TODO delete recipe

	// INGREDIENTS
	h.r.Post("/:id/ingredient", AuthMiddleware(h.db), h.createIngredient)
	h.r.Patch("/:id/ingredient/:ingredientId", AuthMiddleware(h.db), h.updateIngredient)
	h.r.Delete("/:id/ingredient/:ingredientId", AuthMiddleware(h.db), h.deleteIngredient)
	// TODO swap ingredient

	// INSTRUCTIONS
	h.r.Post("/:id/instruction", AuthMiddleware(h.db), h.createInstruction)
	h.r.Patch("/:id/instruction/:instructionId", AuthMiddleware(h.db), h.updateInstruction)
	h.r.Patch("/:id/instruction/:instructionOneId/:instructionTwoId", AuthMiddleware(h.db), h.swapInstructions)
	//	TODO delete ingredient
}

// RECIPES

// POST /recipe/
func (h *RecipeHandler) createRecipe(c *fiber.Ctx) error {
	var user domain.User
	if u, ok := c.Locals("user").(*domain.User); !ok {
		log.Println("no user")
		return SendError(c, Unauthorized())
	} else {
		user = *u
	}

	recipe := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		UserID      uint   `json:"user_id"`
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

	err := h.recipeService.CreateRecipe(recipe.UserID, recipe.Name, recipe.Description)
	if err != nil {
		log.Println("error creating recipe", err)
		return SendError(c, InternalServerError())
	}

	c.Status(fiber.StatusCreated)
	return nil
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
		return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
	}

	return c.JSON(recipe.ToDto())

	//recipe, err := h.recipeService.GetRecipeById(id)
	//if err
}

// PATCH /recipe/:id
func (h *RecipeHandler) updateRecipe(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.UpdateRecipe(uint(id), r.Name, r.Description)
	if err != nil {
		log.Println("error updating recipe", err)
		//TODO handle not found error
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:id
func (h *RecipeHandler) deleteRecipe(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return SendError(c, BadRequest("id must be an integer"))
	}

	err = h.recipeService.DeleteRecipe(uint(id))
	if err != nil {
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// INGREDIENTS

// POST /recipe/:id/ingredient
func (h *RecipeHandler) createIngredient(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.AddIngredientToRecipe(uint(id), ingredient.Name, ingredient.Quantity, ingredient.Unit)
	if err != nil {
		if err.(services.RecipeServiceError).Code == services.ErrRecipeNotFound {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		log.Println("error creating ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/ingredient/:ingredientId
func (h *RecipeHandler) updateIngredient(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.UpdateIngredient(
		uint(recipeID),
		uint(ingredientID),
		ingredient.Name,
		ingredient.Quantity,
		ingredient.Unit,
	)
	if err != nil {
		log.Println("error updating ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// DELETE /recipe/:id/ingredient/:ingredientId
func (h *RecipeHandler) deleteIngredient(c *fiber.Ctx) error {
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

	err = h.recipeService.DeleteIngredient(uint(recipeID), uint(ingredientID))
	if err != nil {
		if err.(services.RecipeServiceError).Code == services.ErrRecipeNotFound {
			return SendError(c, NotFound(map[string]string{"error": "recipe not found"}))
		}
		if err.(services.RecipeServiceError).Code == services.ErrIngredientConflict {
			return SendError(c, Conflict(map[string]string{"error": "ingredient does not belong to recipe"}))
		}
		log.Println("error deleting ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// INSTRUCTIONS

// POST /recipe
func (h *RecipeHandler) createInstruction(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.AddInstructionToRecipe(uint(id), instruction.Step, instruction.Contents)
	if err != nil {
		log.Println("error creating instruction", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/instruction/:instructionId
func (h *RecipeHandler) updateInstruction(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.UpdateInstruction(uint(recipeID), uint(instructionID), instruction.Contents)
	if err != nil {
		log.Println("error updating instruction", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

// PATCH /recipe/:id/instruction/:instructionOneId/:instructionTwoId
func (h *RecipeHandler) swapInstructions(c *fiber.Ctx) error {
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

	recipe, err := h.recipeService.SwapInstructions(uint(recipeID), uint(instructionOneID), uint(instructionTwoID))
	if err != nil {
		log.Println("error swapping instructions", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}
