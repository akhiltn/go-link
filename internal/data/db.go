// Package data provides PostgreSQL-backed persistence for short URLs.
package data

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("short URL not found")

// Driver identifies a database implementation.
type Driver string

const (
	// Postgres selects the PostgreSQL database implementation.
	Postgres Driver = "postgres"
	// Valkey selects the Valkey database implementation.
	Valkey Driver = "valkey"
)

// DB combines short URL persistence operations with database lifecycle management.
type DB interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Delete(context.Context, string) error
	GetAllKeyValues(context.Context) (map[string]string, error)
	Close() error
}

// New creates a database implementation for the selected driver.
func New(ctx context.Context, driver Driver, databaseURL string) (DB, error) {
	switch driver {
	case Postgres:
		return NewPostgresDB(ctx, databaseURL)
	case Valkey:
		return NewValkeyDB(ctx, databaseURL)
	default:
		return nil, fmt.Errorf("unsupported database driver %q", driver)
	}
}
