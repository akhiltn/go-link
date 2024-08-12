package api

import (
	"log"

	"github.com/akhiltn/go-quick-url/internal/data"
	"github.com/akhiltn/go-quick-url/internal/helper"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	Url   string `json:"url"`
	Short string `json:"short"`
}

func GoToURL(c *fiber.Ctx) error {
	key := c.Params("key")

	db, err := data.GetDB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}

	value, err := db.Get(key)
	if err != nil {
		value = "https://www.google.com/search?q=" + key
	}

	log.Printf("Redirecting to: %s", value)
	return c.Redirect(value, fiber.StatusMovedPermanently)
}

func CreateShortURL(c *fiber.Ctx) error {
	db, err := data.GetDB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid request body",
		})
	}
	if !govalidator.IsURL(body.Url) || !helper.RemoveDomainError(body.Url) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid URL",
		})
	}
	body.Url = helper.EnforceHTTP(body.Url)
	err = db.Set(body.Short, body.Url)
	if err != nil {
		log.Printf("Error setting value for key %s: %v", body.Short, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func DeleteShortURL(c *fiber.Ctx) error {
	db, err := data.GetDB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	key := c.Params("key")
	err = db.Delete(key)
	if err != nil {
		log.Printf("Error deleting key %s: %v", key, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetAllKV(c *fiber.Ctx) error {
	db, err := data.GetDB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	kvs, err := db.GetAllKeyValues()
	if err != nil {
		log.Printf("Error getting key-value pairs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	return c.JSON(kvs)
}
