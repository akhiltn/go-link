package api

import (
	"errors"
	"log"
	"strings"

	"github.com/akhiltn/go-link/internal/data"
	"github.com/akhiltn/go-link/internal/helper"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// reservedKeys are served by fixed routes (/api/*, /swagger/*) that must never
// be shadowed by, or shadow, a user's short key.
var reservedKeys = []string{"api", "swagger"}

// isReservedKey is case-insensitive because Fiber routing is too.
func isReservedKey(key string) bool {
	for _, r := range reservedKeys {
		if strings.EqualFold(key, r) {
			return true
		}
	}
	return false
}

type request struct {
	Url   string `json:"url"`
	Short string `json:"short"`
}

// wantsJSON is true only when the client explicitly asks for JSON. Browsers
// send "*/*" on navigation, which must keep redirecting.
func wantsJSON(c *fiber.Ctx) bool {
	return strings.Contains(c.Get(fiber.HeaderAccept), "application/json")
}

// Handler groups the API handlers with their storage dependency.
type Handler struct {
	store  data.DB
	domain string
}

func NewHandler(store data.DB, domain string) *Handler {
	return &Handler{store: store, domain: domain}
}

// ResolveShortURL method to redirect to the original URL.
// @Description Redirect to the original URL (302). With "Accept: application/json" (e.g. Swagger "Try it out") the target is returned as JSON instead of redirecting.
// @Summary Redirect to the original URL.
// @Tags api
// @Produce json
// @Param key path string true "Key"
// @Success 200 {object} map[string]string "Key and target URL (JSON clients)"
// @Success 302 {string} string "Found"
// @Router /{key} [get]
func (h *Handler) ResolveShortURL(c *fiber.Ctx) error {
	key := c.Params("key")
	value, err := h.store.Get(c.UserContext(), key)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			value = "https://www.google.com/search?q=" + key
		} else {
			log.Printf("Error getting key %s: %v", key, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "Internal server error",
			})
		}
	}
	if wantsJSON(c) {
		return c.JSON(map[string]string{key: value})
	}
	log.Printf("Redirecting to: %s", value)
	return c.Redirect(value, fiber.StatusFound)
}

// CreateShortURL method to create a short URL.
// @Description Create a short URL.
// @Summary Create a short URL.
// @Tags api
// @Accept json
// @Param url body request true "URL"
// @Success 201 {string} string "Created"
// @Router / [post]
func (h *Handler) CreateShortURL(c *fiber.Ctx) error {
	var body request
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid request body",
		})
	}
	if !govalidator.IsURL(body.Url) || !helper.RemoveDomainError(body.Url, h.domain) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid URL",
		})
	}
	body.Url = helper.EnforceHTTP(body.Url)
	if isReservedKey(body.Short) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Short key is reserved",
		})
	}
	err := h.store.Set(c.UserContext(), body.Short, body.Url)
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
func (h *Handler) DeleteShortURL(c *fiber.Ctx) error {
	key := c.Params("key")
	err := h.store.Delete(c.UserContext(), key)
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
// @Router /api/allkv [get]
func (h *Handler) GetAllKV(c *fiber.Ctx) error {
	kvs, err := h.store.GetAllKeyValues(c.UserContext())
	if err != nil {
		log.Printf("Error getting key-value pairs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}
	return c.JSON(kvs)
}
