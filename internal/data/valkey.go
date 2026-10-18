package data

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

const valkeyKeyPrefix = "go-link:"

type valkeyDB struct {
	client *redis.Client
}

var _ DB = (*valkeyDB)(nil)

func (v *valkeyDB) Close() error {
	return v.client.Close()
}

// NewValkeyDB opens a Valkey connection at address.
func NewValkeyDB(ctx context.Context, address string) (DB, error) {
	if address == "" {
		return nil, errors.New("valkey address is required")
	}

	client := redis.NewClient(&redis.Options{Addr: address})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &valkeyDB{client: client}, nil
}

func (v *valkeyDB) Get(ctx context.Context, key string) (string, error) {
	value, err := v.client.Get(ctx, valkeyKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (v *valkeyDB) Set(ctx context.Context, key string, value string) error {
	if key == "" {
		return errors.New("short key is required")
	}
	return v.client.Set(ctx, valkeyKey(key), value, 0).Err()
}

func (v *valkeyDB) Delete(ctx context.Context, key string) error {
	return v.client.Del(ctx, valkeyKey(key)).Err()
}

func (v *valkeyDB) GetAllKeyValues(ctx context.Context) (map[string]string, error) {
	values := make(map[string]string)
	var cursor uint64
	for {
		keys, nextCursor, err := v.client.Scan(ctx, cursor, valkeyKeyPrefix+"*", 0).Result()
		if err != nil {
			return nil, err
		}
		for _, storageKey := range keys {
			value, err := v.client.Get(ctx, storageKey).Result()
			if errors.Is(err, redis.Nil) {
				continue
			}
			if err != nil {
				return nil, err
			}
			values[strings.TrimPrefix(storageKey, valkeyKeyPrefix)] = value
		}
		if nextCursor == 0 {
			return values, nil
		}
		cursor = nextCursor
	}
}

func valkeyKey(key string) string {
	return fmt.Sprintf("%s%s", valkeyKeyPrefix, key)
}
