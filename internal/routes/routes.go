package routes

import (
	"github.com/akhiltn/go-link/internal/api"
	"github.com/gofiber/fiber/v2"
)

func AppRouteInit(app *fiber.App, store api.Store) {
	handler := api.NewHandler(store)
	app.Get("/allkv", handler.GetAllKV)
	app.Get("/:key", handler.ResolveShortURL)
	app.Delete("/:key", handler.DeleteShortURL)
	app.Post("/", handler.CreateShortURL)
}
