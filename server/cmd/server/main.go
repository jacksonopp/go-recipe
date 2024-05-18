package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/handlers"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	db, err := createDb()
	if err != nil {
		log.Panicf("failed to create database %v", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	api := app.Group("/api")

	userHandler := handlers.NewUserHandler(api, db)
	recipeHandler := handlers.NewRecipeHandler(api, db)
	createApiRoutes(userHandler, recipeHandler)

	sessionService := services.NewSessionService(db)
	// TODO figure out what to do with done channel
	done, err := sessionService.PruneOnSchedule(time.Hour * 24)
	if err != nil {
		done <- true
		log.Printf("ERROR: failed to prune session %v", err)
	}

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

func createDb() (*gorm.DB, error) {
	dsn := createDsn()
	log.Println("connecting to database", dsn)

	//dsn := "host=localhost user=postgres password=postgres dbname=go_recipe port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Session{},
		&domain.Recipe{},
		&domain.Ingredient{},
		&domain.Instruction{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createApiRoutes(handlers ...handlers.Handler) {
	for _, handler := range handlers {
		handler.RegisterRoutes()
	}
}
