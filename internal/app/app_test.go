package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akhiltn/go-link/internal/data"
	"github.com/gofiber/fiber/v2"
)

type memoryRepository struct {
	values map[string]string
}

func (r *memoryRepository) Get(_ context.Context, key string) (string, error) {
	value, ok := r.values[key]
	if !ok {
		return "", data.ErrNotFound
	}
	return value, nil
}

func (r *memoryRepository) Set(_ context.Context, key, value string) error {
	r.values[key] = value
	return nil
}

func (r *memoryRepository) Delete(_ context.Context, key string) error {
	delete(r.values, key)
	return nil
}

func (r *memoryRepository) GetAllKeyValues(context.Context) (map[string]string, error) {
	return r.values, nil
}

func (r *memoryRepository) Close() error { return nil }

func TestNewPassesConfiguredDomainToCreateRoute(t *testing.T) {
	repository := &memoryRepository{values: make(map[string]string)}
	server := New(repository, Config{Domain: "localhost"})

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{"short":"custom","url":"http://localhost:8080/path"}`,
	))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := server.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("POST same-domain status = %d, want %d", response.StatusCode, fiber.StatusBadRequest)
	}
	if _, ok := repository.values["custom"]; ok {
		t.Fatal("same-domain URL was stored")
	}

	request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{"short":"external","url":"https://example.org/path"}`,
	))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err = server.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusCreated {
		t.Fatalf("POST external status = %d, want %d", response.StatusCode, fiber.StatusCreated)
	}
	if got := repository.values["external"]; got != "https://example.org/path" {
		t.Fatalf("stored URL = %q, want %q", got, "https://example.org/path")
	}
}
