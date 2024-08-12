package routes

import (
	"github.com/akhiltn/go-quick-url/internal/api"
	"github.com/gofiber/fiber/v2"
)

func AppRouteInit(app *fiber.App) {
	app.Get("/:key", api.GoToURL)
	app.Post("/", api.CreateShortURL)
}
