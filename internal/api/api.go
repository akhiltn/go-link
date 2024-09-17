package api

import (
	"log"

	"github.com/akhiltn/go-link/internal/data"
	"github.com/akhiltn/go-link/internal/helper"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	Url   string `json:"url"`
	Short string `json:"short"`
}

// ResolveShortURL method to redirect to the original URL.
// @Description Redirect to the original URL.
// @Summary Redirect to the original URL.
// @Tags api
// @Produce json
// @Param key path string true "Key"
// @Success 301 {string} string "Moved Permanently"
// @Router /{key} [get]
func ResolveShortURL(c *fiber.Ctx) error {
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

// CreateShortURL method to create a short URL.
// @Description Create a short URL.
// @Summary Create a short URL.
// @Tags api
// @Accept json
// @Param url body request true "URL"
// @Success 201 {string} string "Created"
// @Router / [post]
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
	log.Printf("Created key %s with value %s", body.Short, body.Url)
	return c.SendStatus(fiber.StatusCreated)
}

// DeleteShortURL method to delete a short URL.
// @Description Delete a short URL.
// @Summary Delete a short URL.
// @Tags api
// @Produce json
// @Param key path string true "Key"
// @Success 200 {string} string "OK"
// @Router /{key} [delete]
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

// GetAllKV method to get all key-value pairs.
// @Description Get all key-value pairs.
// @Summary Get all key-value pairs.
// @Tags api
// @Produce json
// @Success 200 {string} string "OK"
// @Router /allkv [get]
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
