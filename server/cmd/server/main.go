package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=recipe port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
		return
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Panicf("failed to migrate database: %v", err)
		return
	}

	app := fiber.New()
	api := app.Group("/api")

	userHandler := handlers.NewUserHandler(api, db)
	createApiRoutes(userHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Listen(":8080")
}

func createApiRoutes(handlers ...handlers.Handler) {
	for _, handler := range handlers {
		handler.CreateAllRoutes()
	}
}
