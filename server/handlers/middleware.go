package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from the header
		sessionCookie := struct {
			Session string `cookie:"session"`
		}{}

		if err := c.CookieParser(&sessionCookie); err != nil {
			return SendError(c, Unauthorized())
		}

		token := sessionCookie.Session

		if token == "" {
			return SendError(c, Unauthorized())
		}

		user, err := getUserBySessionToken(db, token)
		if err != nil {
			return SendError(c, Unauthorized())
		}

		// Pass the user to the next handler
		c.Locals("user", user)
		return c.Next()
	}
}

func getUserBySessionToken(db *gorm.DB, token string) (*domain.User, error) {
	var user domain.User
	err := db.Table("users").
		Joins("inner join sessions on users.id = sessions.user_id").
		Where("sessions.token = ?", token).
		Where("sessions.deleted_at is null").
		First(&user).
		Error
	return &user, err
}
