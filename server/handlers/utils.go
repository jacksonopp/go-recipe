package handlers

import (
	"github.com/gofiber/fiber/v2"
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
