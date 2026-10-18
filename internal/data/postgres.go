package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/akhiltn/go-link/internal/data/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const postgresMigrationLockID int64 = 0x676f6c696e6b

type postgresDB struct {
	db *sql.DB
}

var _ DB = (*postgresDB)(nil)

var (
	gooseConfigOnce sync.Once
	gooseConfigErr  error
)

func (p *postgresDB) Close() error {
	return p.db.Close()
}

// NewPostgresDB opens a PostgreSQL database and applies pending migrations.
func NewPostgresDB(ctx context.Context, databaseURL string) (DB, error) {
	if databaseURL == "" {
		return nil, errors.New("postgres database URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := configureGoose(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrateWithLock(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &postgresDB{db: db}, nil
}

func configureGoose() error {
	gooseConfigOnce.Do(func() {
		goose.SetBaseFS(migrations.FS())
		gooseConfigErr = goose.SetDialect("postgres")
	})
	return gooseConfigErr
}

func migrateWithLock(ctx context.Context, db *sql.DB) (err error) {
	lockConn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("get migration lock connection: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, unlockErr := lockConn.ExecContext(unlockCtx, "SELECT pg_advisory_unlock($1)", postgresMigrationLockID); unlockErr != nil {
			_ = lockConn.Raw(func(any) error { return driver.ErrBadConn })
			err = errors.Join(err, fmt.Errorf("release migration lock: %w", unlockErr))
		}
		if closeErr := lockConn.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migration lock connection: %w", closeErr))
		}
	}()

	if _, err := lockConn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", postgresMigrationLockID); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("apply embedded migrations: %w", err)
	}
	return nil
}

func (p *postgresDB) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := p.db.QueryRowContext(ctx, "SELECT target_url FROM short_url_mappings WHERE short_key = $1", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (p *postgresDB) Set(ctx context.Context, key string, value string) error {
	if key == "" {
		return errors.New("short key is required")
	}
	_, err := p.db.ExecContext(ctx, `INSERT INTO short_url_mappings (short_key, target_url) VALUES ($1, $2)
ON CONFLICT (short_key) DO UPDATE SET target_url = EXCLUDED.target_url`, key, value)
	return err
}

func (p *postgresDB) Delete(ctx context.Context, key string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM short_url_mappings WHERE short_key = $1", key)
	return err
}

func (p *postgresDB) GetAllKeyValues(ctx context.Context) (map[string]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT short_key, target_url FROM short_url_mappings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		values[key] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
