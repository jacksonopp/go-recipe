package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	dsn := createDsn()
	log.Println("connecting to database", dsn)

	//dsn := "host=localhost user=postgres password=postgres dbname=go_recipe port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
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
	app.Use(logger.New())
	api := app.Group("/api")

	userHandler := handlers.NewUserHandler(api, db)
	createApiRoutes(userHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("request to /")

		return c.SendString("ok")
	})

	err = app.Listen("0.0.0.0:8080")
	if err != nil {
		log.Panicf("failed to start server: %v", err)
	}
}

func createDsn() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
}

func createApiRoutes(handlers ...handlers.Handler) {
	for _, handler := range handlers {
		handler.CreateAllRoutes()
	}
}
