package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService services.UserService
	r           fiber.Router
	db          *gorm.DB
}

func NewUserHandler(r fiber.Router, db *gorm.DB) *UserHandler {
	subpath := r.Group("/user")
	userService := services.NewUserService(db)
	return &UserHandler{userService: userService, r: subpath, db: db}
}

func (h *UserHandler) RegisterRoutes() {
	h.r.Get("/:name", h.getUserByName)
	h.r.Get("/:name/recipes", h.getUserRecipes)
}

func (h *UserHandler) getUserByName(c *fiber.Ctx) error {
	username := c.Params("name")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("username is required")
	}
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(user.ToDto())
}

// GET /user/:name/recipes?page={n}&limit={n}
func (h *UserHandler) getUserRecipes(c *fiber.Ctx) error {
	page, limit := getPaginationParams(c)
	username := c.Params("name")
	if username == "" {
		return SendError(c, BadRequest("username is required"))
	}

	recipes, err := h.userService.GetUsersRecipes(username, page, limit)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return SendError(c, NotFound(map[string]string{"msg": "user not found"}))
		}
		return SendError(c, InternalServerError())
	}
	recipesDtos := make([]domain.RecipeDto, len(recipes))
	for i, r := range recipes {
		rdto, ok := r.ToDto().(domain.RecipeDto)
		if !ok {
			return SendError(c, InternalServerError())
		}
		recipesDtos[i] = rdto
	}
	return c.JSON(recipesDtos)
}
