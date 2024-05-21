package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
	"log"
)

type AuthHandler struct {
	r              fiber.Router
	authService    services.AuthService
	sessionService services.SessionService
}

func NewAuthHandler(r fiber.Router, db *gorm.DB) *AuthHandler {
	authService := services.NewAuthService(db)
	sessionService := services.NewSessionService(db)

	subpath := r.Group("/auth")

	return &AuthHandler{r: subpath, authService: authService, sessionService: sessionService}
}

func (h *AuthHandler) RegisterRoutes() {
	h.r.Post("/register", h.register)
	h.r.Post("/login", h.login)
	h.r.Get("/session", h.session)
}

func (h *AuthHandler) register(c *fiber.Ctx) error {
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

	if err := h.authService.CreateUser(u); err != nil {
		var e services.AuthServiceError
		if errors.As(err, &e) {
			if e.Code == services.ErrUserAlreadyExists {
				err := Conflict(map[string]string{"username": e.Msg})
				return SendError(c, err)
			}
		}

		log.Printf("error creating user: %v", err)
		return InternalServerError()
	}
	return c.Redirect("/login", fiber.StatusPermanentRedirect)
}

func (h *AuthHandler) login(c *fiber.Ctx) error {
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

	u, err := h.authService.LoginUser(user.Username, user.Password)
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

	return c.JSON(u.ToDto())
}

func (h *AuthHandler) session(c *fiber.Ctx) error {
	session := c.Cookies("session")
	if session == "" {
		err := Unauthorized()
		return SendError(c, err)
	}

	err := h.sessionService.CheckSession(session)
	if err != nil {
		log.Println("error getting user by session: ", err)
		err := Unauthorized()
		return SendError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
