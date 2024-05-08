package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
	"log"
)

type UserHandler struct {
	// userService services.UserService
	r fiber.Router
	s *services.UserService
}

func NewUserHandler(r fiber.Router, db *gorm.DB) *UserHandler {
	s := services.NewUserService(db)

	return &UserHandler{r: r, s: s}
}

func (h *UserHandler) CreateAllRoutes() {
	h.r.Post("/login", h.login)
	// h.r.Post("/register", h.register)
}

func (h *UserHandler) login(c *fiber.Ctx) error {
	user := struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}{}

	if err := c.BodyParser(&user); err != nil {
		log.Printf("error parsing body: %v", err)
		return c.Status(400).SendString("bad request")
	}

	if user.Username == "" {
		err := UnprocessableEntity(map[string]string{"username": "username is required"})
		return SendError(c, err)
	}

	if user.Password == "" {
		err := UnprocessableEntity(map[string]string{"password": "password is required"})
		return SendError(c, err)
	}

	if user.Password != user.PasswordConfirm {
		err := UnprocessableEntity(map[string]string{"password": "passwords do not match"})
		return SendError(c, err)
	}

	u := domain.User{
		Username: user.Username,
		Password: user.Password,
	}

	if err := h.s.CreateUser(u); err != nil {
		log.Printf("error creating user: %v", err)
		return SendError(c, InternalServerError())
	}

	return c.SendString("login")
}
