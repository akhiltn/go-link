package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akhiltn/go-link/internal/data"
	"github.com/gofiber/fiber/v2"
)

type memStore struct {
	m      map[string]string
	getErr error
	setErr func(string, string) error
}

func (s *memStore) Get(_ context.Context, k string) (string, error) {
	if s.getErr != nil {
		return "", s.getErr
	}
	return s.m[k], nil
}
func (s *memStore) Set(_ context.Context, k, v string) error {
	if s.setErr != nil {
		if err := s.setErr(k, v); err != nil {
			return err
		}
	}
	s.m[k] = v
	return nil
}
func (s *memStore) Delete(_ context.Context, k string) error { delete(s.m, k); return nil }
func (s *memStore) GetAllKeyValues(context.Context) (map[string]string, error) {
	return s.m, nil
}
func (s *memStore) Close() error { return nil }

func TestCreateShortURLRejectsConfiguredDomain(t *testing.T) {
	store := &memStore{m: map[string]string{}}
	app := fiber.New()
	app.Post("/", NewHandler(store, "example.com").CreateShortURL)

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"short":"test","url":"example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("got %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}
	if _, ok := store.m["test"]; ok {
		t.Fatal("same-domain URL was stored")
	}
}

func TestCreateShortURLRejectsNullBody(t *testing.T) {
	app := fiber.New()
	app.Post("/", NewHandler(&memStore{m: map[string]string{}}, "example.com").CreateShortURL)

	req := httptest.NewRequest("POST", "/", strings.NewReader("null"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("got %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}
}

func TestCreateShortURLReservedKeys(t *testing.T) {
	store := &memStore{m: map[string]string{}}
	app := fiber.New()
	app.Post("/", NewHandler(store, "example.com").CreateShortURL)

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
		body := `{"short":"` + tc.short + `","url":"example.org"}`
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

func TestCreateShortURLStorageFailureForEmptyShort(t *testing.T) {
	store := &memStore{
		m: map[string]string{},
		setErr: func(k, _ string) error {
			if k == "" {
				return errors.New("short key is required")
			}
			return nil
		},
	}
	app := fiber.New()
	app.Post("/", NewHandler(store, "example.com").CreateShortURL)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "omitted short",
			body: `{"url":"https://example.org/path"}`,
		},
		{
			name: "explicit empty short",
			body: `{"short":"","url":"https://example.org/path"}`,
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("%s: app.Test() error = %v", tc.name, err)
		}
		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Fatalf("%s: got %d, want %d", tc.name, resp.StatusCode, fiber.StatusInternalServerError)
		}
		if _, ok := store.m[""]; ok {
			t.Fatalf("%s: empty key was stored", tc.name)
		}
	}
}

func TestResolveShortURL(t *testing.T) {
	store := &memStore{m: map[string]string{"a.b": "https://www.example.org"}}
	app := fiber.New()
	app.Get("/:key", NewHandler(store, "example.com").ResolveShortURL)

	tests := []struct {
		name, accept string
		want         int
		wantLoc      string
	}{
		{"browser redirects", "text/html,application/xhtml+xml,*/*;q=0.8", fiber.StatusFound, "https://www.example.org"},
		{"no accept redirects", "", fiber.StatusFound, "https://www.example.org"},
		{"json client gets body", "application/json", fiber.StatusOK, ""},
	}
	for _, tc := range tests {
		req := httptest.NewRequest("GET", "/a.b", nil)
		if tc.accept != "" {
			req.Header.Set("Accept", tc.accept)
		}
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, resp.StatusCode, tc.want)
		}
		if loc := resp.Header.Get("Location"); loc != tc.wantLoc {
			t.Errorf("%s: Location %q, want %q", tc.name, loc, tc.wantLoc)
		}
		if tc.want == fiber.StatusOK {
			var got map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if got["a.b"] != "https://www.example.org" || len(got) != 1 {
				t.Errorf("%s: body %v", tc.name, got)
			}
		}
	}
}

func TestResolveShortURLMissingLinkFallsBackToSearch(t *testing.T) {
	store := &memStore{m: map[string]string{}, getErr: data.ErrNotFound}
	app := fiber.New()
	app.Get("/:key", NewHandler(store, "example.com").ResolveShortURL)

	req := httptest.NewRequest("GET", "/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusFound {
		t.Fatalf("got %d, want %d", resp.StatusCode, fiber.StatusFound)
	}
	if got := resp.Header.Get("Location"); got != "https://www.google.com/search?q=abc" {
		t.Fatalf("got Location %q, want %q", got, "https://www.google.com/search?q=abc")
	}
}

func TestResolveShortURLStorageErrorReturnsServerError(t *testing.T) {
	store := &memStore{m: map[string]string{}, getErr: errors.New("database unavailable")}
	app := fiber.New()
	app.Get("/:key", NewHandler(store, "example.com").ResolveShortURL)

	req := httptest.NewRequest("GET", "/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("got %d, want %d", resp.StatusCode, fiber.StatusInternalServerError)
	}
	if resp.Header.Get("Location") != "" {
		t.Fatalf("unexpected redirect to %q", resp.Header.Get("Location"))
	}
}
