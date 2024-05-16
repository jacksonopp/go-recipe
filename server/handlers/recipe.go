package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type RecipeHandler struct {
	r  fiber.Router
	db *gorm.DB
}

func NewRecipeHandler(r fiber.Router, db *gorm.DB) *RecipeHandler {
	subpath := r.Group("/recipe")

	return &RecipeHandler{r: subpath, db: db}
}

func (h *RecipeHandler) CreateAllRoutes() {
	h.r.Get("/demo", AuthMiddleware(h.db), h.demo)
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
