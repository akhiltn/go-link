package app

import (
	"github.com/akhiltn/go-link/internal/data"
	"github.com/akhiltn/go-link/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/swagger"
)

type Config struct {
	Domain string
}

func New(db data.DB, config Config) *fiber.App {
	server := fiber.New()
	server.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:8080",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	server.Get("/swagger/*", swagger.HandlerDefault)
	server.Use(healthcheck.New())
	routes.AppRouteInit(server, db, config.Domain)
	return server
}
