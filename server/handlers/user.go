package handlers

import (
	"github.com/gofiber/fiber/v2"
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

func (h *UserHandler) getUserRecipes(c *fiber.Ctx) error {
	return c.SendString("ok")
}
