package main

import (
	"log"
	"os"

	"github.com/akhiltn/go-quick-url/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	app.Use(healthcheck.New())
  routes.AppRouteInit(app)
	log.Printf("Server started on port %s", os.Getenv("APP_PORT"))
	app.Listen(os.Getenv("APP_PORT"))
}
