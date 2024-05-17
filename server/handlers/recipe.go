package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
	"log"
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
	h.r.Get("/demo", AuthMiddleware(h.db), h.demo)
	h.r.Post("/", AuthMiddleware(h.db), h.createRecipe)
}

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

func (h *RecipeHandler) demo(c *fiber.Ctx) error {
	var user domain.User
	if u, ok := c.Locals("user").(*domain.User); !ok {
		log.Println("no user")
		return SendError(c, Unauthorized())
	} else {
		user = *u
	}

	return c.JSON(user.ToDto())
}
