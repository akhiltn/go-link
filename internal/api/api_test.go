package api

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type memStore struct{ m map[string]string }

func (s *memStore) Get(_ context.Context, k string) (string, error) { return s.m[k], nil }
func (s *memStore) Set(_ context.Context, k, v string) error        { s.m[k] = v; return nil }
func (s *memStore) Delete(_ context.Context, k string) error        { delete(s.m, k); return nil }
func (s *memStore) GetAllKeyValues(context.Context) (map[string]string, error) {
	return s.m, nil
}

func TestCreateShortURLReservedKeys(t *testing.T) {
	store := &memStore{m: map[string]string{}}
	app := fiber.New()
	app.Post("/", NewHandler(store).CreateShortURL)

	tests := []struct {
		short string
		want  int
	}{
		{"api", fiber.StatusBadRequest},
		{"API", fiber.StatusBadRequest},
		{"swagger", fiber.StatusBadRequest},
		{"Swagger", fiber.StatusBadRequest},
		{"allkv", fiber.StatusCreated},
		{"test", fiber.StatusCreated},
	}
	for _, tc := range tests {
		body := `{"short":"` + tc.short + `","url":"example.com"}`
		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != tc.want {
			t.Errorf("short %q: got %d, want %d", tc.short, resp.StatusCode, tc.want)
		}
	}
	for _, k := range []string{"api", "API", "swagger", "Swagger"} {
		if _, ok := store.m[k]; ok {
			t.Errorf("reserved key %q was stored", k)
		}
	}
}
