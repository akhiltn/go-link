package api

import (
	"github.com/akhiltn/go-quick-url/internal/data"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	Url   string `json:"url"`
	Short string `json:"short"`
}

func RedirectToURL(c *fiber.Ctx) error {
	key := c.Params("key")
	db, err := data.InitDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer db.Close()
	value, err := db.Get(c.Params(key))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.Redirect(value, 301)
}

func CreateShortURL(c *fiber.Ctx) error {
  db, err := data.InitDB()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": true,
      "msg":   err.Error(),
    })
  }
  defer db.Close()
  body := new(request)
  if err := c.BodyParser(&body); err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": true,
      "msg":   err.Error(),
    })
  }
  err = db.Set(body.Short, body.Url)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": true,
      "msg":   err.Error(),
    })
  }
  return c.SendStatus(fiber.StatusCreated)
}
