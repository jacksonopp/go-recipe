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

func (h *RecipeHandler) CreateAllRoutes() {
	h.r.Post("/", AuthMiddleware(h.db), h.createRecipe)
	h.r.Get("/:id", h.getRecipeById)
	h.r.Post("/:id/ingredient", AuthMiddleware(h.db), h.createIngredient)
	h.r.Post("/:id/instruction", AuthMiddleware(h.db), h.createInstruction)
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

// POST /recipe/:id/ingredient
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

	err := h.recipeService.CreateRecipe(recipe.Name, recipe.Description, recipe.UserID)
	if err != nil {
		log.Println("error creating recipe", err)
		return SendError(c, InternalServerError())
	}

	c.Status(fiber.StatusCreated)
	return nil
}

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
		log.Println("error creating ingredient", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(recipe.ToDto())
}

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
