package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewValkeyDBRejectsEmptyAddress(t *testing.T) {
	_, err := NewValkeyDB(context.Background(), "")
	if err == nil {
		t.Fatal("NewValkeyDB() error = nil, want error for an empty address")
	}
}

func TestValkeyRepository(t *testing.T) {
	address := os.Getenv("TEST_VALKEY_ADDR")
	if address == "" {
		t.Skip("set TEST_VALKEY_ADDR to run Valkey repository tests")
	}

	ctx := context.Background()
	repo, err := New(ctx, Valkey, address)
	if err != nil {
		t.Fatal(err)
	}
	key := fmt.Sprintf("test-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_ = repo.Delete(context.Background(), key)
		_ = repo.Close()
	})

	if err := repo.Set(ctx, key, "https://example.com/first"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Set(ctx, key, "https://example.com/updated"); err != nil {
		t.Fatal(err)
	}
	got, err := repo.Get(ctx, key)
	if err != nil || got != "https://example.com/updated" {
		t.Fatalf("Get() = %q, %v; want updated URL, nil", got, err)
	}
	all, err := repo.GetAllKeyValues(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if all[key] != "https://example.com/updated" {
		t.Fatalf("GetAllKeyValues()[%q] = %q; want updated URL", key, all[key])
	}
	if _, err := repo.Get(ctx, key+"-missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing Get() error = %v; want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, key); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, key); err != nil {
		t.Fatalf("second Delete() = %v; want nil", err)
	}
	all, err = repo.GetAllKeyValues(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := all[key]; ok {
		t.Fatalf("deleted key %q is still listed", key)
	}
}
