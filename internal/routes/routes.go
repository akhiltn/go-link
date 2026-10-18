package routes

import (
	"github.com/akhiltn/go-link/internal/api"
	"github.com/akhiltn/go-link/internal/data"
	"github.com/gofiber/fiber/v2"
)

func AppRouteInit(app *fiber.App, db data.DB, domain string) {
	handler := api.NewHandler(db, domain)
	app.Get("/api/allkv", handler.GetAllKV)
	app.Get("/:key", handler.ResolveShortURL)
	app.Delete("/:key", handler.DeleteShortURL)
	app.Post("/", handler.CreateShortURL)
}
