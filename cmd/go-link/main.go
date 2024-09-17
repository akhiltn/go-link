package main

import (
	"log"
	"os"

	_ "github.com/akhiltn/go-link/docs"
	"github.com/akhiltn/go-link/internal/data"
	"github.com/akhiltn/go-link/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Go Link
// @version 1.0
// @description URL shortener app
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3100
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	db, err := data.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := fiber.New()
	// Configure CORS settings
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:8080", // Allow localhost origins
		AllowMethods: "GET, POST, PUT, DELETE",                       // Allowed methods
		AllowHeaders: "Origin, Content-Type, Accept",                 // Allowed headers
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
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
