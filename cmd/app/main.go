package main

import (
	"log"
	"os"

	"github.com/akhiltn/go-quick-url/internal/data"
	"github.com/akhiltn/go-quick-url/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := data.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()
	app.Use(healthcheck.New())
	routes.AppRouteInit(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port if not set
	}
	log.Printf("Server started on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
