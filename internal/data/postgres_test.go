package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

const testMigrationLockID int64 = 0x676f6c696e6b

func TestDBSetRejectsEmptyKey(t *testing.T) {
	repo := &postgresDB{}

	err := repo.Set(context.Background(), "", "https://example.com")
	if err == nil {
		t.Fatal("Set() error = nil, want non-nil for empty key")
	}
}

func TestPostgresRepository(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL repository tests")
	}
	ctx := context.Background()
	repo, err := NewPostgresDB(ctx, dsn)
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

func TestOpenMigrationLockHonorsContext(t *testing.T) {
	dsn := postgresTestSchema(t)
	admin, err := sql.Open("pgx", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	lockConn, err := admin.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer lockConn.Close()
	if _, err := lockConn.ExecContext(context.Background(), "SELECT pg_advisory_lock($1)", testMigrationLockID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := lockConn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", testMigrationLockID); err != nil {
			t.Errorf("release migration lock: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err = NewPostgresDB(ctx, dsn)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("NewPostgresDB() error = %v; want context deadline while migration lock is held", err)
	}
}

func TestOpenSerializesConcurrentMigrations(t *testing.T) {
	dsn := postgresTestSchema(t)
	start := make(chan struct{})
	results := make(chan struct {
		repo DB
		err  error
	}, 2)

	for range 2 {
		go func() {
			<-start
			repo, err := NewPostgresDB(context.Background(), dsn)
			results <- struct {
				repo DB
				err  error
			}{repo: repo, err: err}
		}()
	}
	close(start)

	var openErrors []error
	for range 2 {
		result := <-results
		if result.repo != nil {
			if err := result.repo.Close(); err != nil {
				t.Errorf("close repository: %v", err)
			}
		}
		if result.err != nil {
			openErrors = append(openErrors, result.err)
		}
	}
	if len(openErrors) != 0 {
		t.Fatalf("concurrent NewPostgresDB() errors: %v", openErrors)
	}
}

func postgresTestSchema(t *testing.T) string {
	t.Helper()
	baseDSN := os.Getenv("TEST_DATABASE_URL")
	if baseDSN == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}

	admin, err := sql.Open("pgx", baseDSN)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := admin.PingContext(ctx); err != nil {
		_ = admin.Close()
		t.Fatal(err)
	}
	schema := fmt.Sprintf("go_link_test_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA "`+schema+`"`); err != nil {
		_ = admin.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA "`+schema+`" CASCADE`); err != nil {
			t.Errorf("drop test schema: %v", err)
		}
		if err := admin.Close(); err != nil {
			t.Errorf("close admin connection: %v", err)
		}
	})

	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
