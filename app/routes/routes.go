package routes

import (
	"github.com/akhiltn/go-quick-url/app/api"
	"github.com/gofiber/fiber/v2"
)

func AppRouteInit(app *fiber.App) {
	app.Get("/allkv", api.GetAllKV)
	app.Get("/:key", api.GoToURL)
	app.Delete("/:key", api.DeleteShortURL)
	app.Post("/", api.CreateShortURL)
}
