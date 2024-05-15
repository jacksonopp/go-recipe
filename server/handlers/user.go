package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
	"log"
)

type UserHandler struct {
	r              fiber.Router
	userService    services.UserService
	sessionService services.SessionService
}

func NewUserHandler(r fiber.Router, db *gorm.DB) *UserHandler {
	userService := services.NewUserService(db)
	sessionService := services.NewSessionService(db)

	return &UserHandler{r: r, userService: userService, sessionService: sessionService}
}

func (h *UserHandler) CreateAllRoutes() {
	h.r.Post("/register", h.register)
	h.r.Post("/login", h.login)
}

func (h *UserHandler) register(c *fiber.Ctx) error {
	user := struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}{}

	if err := c.BodyParser(&user); err != nil {
		log.Printf("body content: %v", string(c.Body()))
		log.Printf("error parsing body: %v", err)
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
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

	if err := h.userService.CreateUser(u); err != nil {
		var e services.UserServiceError
		if errors.As(err, &e) {
			if e.Code == services.ErrUserAlreadyExists {
				err := Conflict(map[string]string{"username": e.Msg})
				return SendError(c, err)
			}
		}

		log.Printf("error creating user: %v", err)
		return InternalServerError()
	}
	c.Status(201)

	return nil
}

func (h *UserHandler) login(c *fiber.Ctx) error {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&user); err != nil {
		err := UnprocessableEntity(map[string]string{"error": "invalid request body"})
		return SendError(c, err)
	}

	if user.Username == "" {
		err := UnprocessableEntity(map[string]string{"username": "username is required"})
		return SendError(c, err)
	}

	if user.Password == "" {
		err := UnprocessableEntity(map[string]string{"password": "password is required"})
		return SendError(c, err)
	}

	u, err := h.userService.LoginUser(user.Username, user.Password)
	if err != nil {
		log.Println("error logging in user: ", err)
		err := NotFound(map[string]string{"error": "user not found"})
		return SendError(c, err)
	}

	token, err := h.sessionService.CreateSession(u.ID)
	if err != nil {
		log.Println("error creating session: ", err)
		return SendError(c, InternalServerError())
	}

	c.Cookie(&fiber.Cookie{
		Name:  "session",
		Value: token,
	})

	return c.JSON(u)
}
