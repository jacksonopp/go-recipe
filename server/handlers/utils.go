package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"log"
	"strconv"
)

func getPaginationParams(c *fiber.Ctx) (int, int) {
	var (
		page  int
		limit int
		err   error
	)

	page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	log.Println(page, limit)
	return page, limit
}

func getUserFromLocals(c *fiber.Ctx) (domain.User, error) {
	var user domain.User
	if u, ok := c.Locals("user").(*domain.User); !ok {
		log.Println("no user")
		return domain.User{}, Unauthorized()
	} else {
		user = *u
	}
	return user, nil
}
