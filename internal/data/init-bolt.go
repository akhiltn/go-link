package data

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boltdb/bolt"
)

const (
	dbname = "DB"
)

func newBoltDB(dbPath string) (*BoltDB, error) {
	db, err := bolt.Open(os.ExpandEnv(dbPath), 0600, nil)
	if err != nil {
		log.Printf("Error opening BoltDB: %v", err)
		return nil, err
	}
	return &BoltDB{db: db}, nil
}

func OpenBoltDB(dbPath string) (*BoltDB, error) {
	log.Println("Initializing BoltDB")
	if dbPath == "" && runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Println("Error getting home directory:", err)
			return nil, err
		}
		dbPath = filepath.Join(homeDir, ".go-link", "go-link.db")
		if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
			if err != nil {
				log.Printf("Failed to create directory: %v", err)
				return nil, err
			}
		}
	} else if dbPath == "" {
		dbPath = "go-link.db"
	}
	instance, err := newBoltDB(dbPath)
	if err != nil {
		log.Printf("Failed to initialize BoltDB: %v", err)
		return nil, err
	} else {
		// Ensure the bucket is created
		err = instance.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dbname))
			return err
		})
		if err != nil {
			log.Printf("Failed to create bucket: %v", err)
			_ = instance.Close()
			return nil, err
		}
	}
	return instance, err
}

func (b *BoltDB) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var value string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbname))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		val := bucket.Get([]byte(key))
		if val == nil {
			return bolt.ErrBucketNotFound
		}
		value = string(val)
		return nil
	})
	return value, err
}

func (b *BoltDB) Set(ctx context.Context, key string, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(dbname))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (b *BoltDB) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbname))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func (b *BoltDB) GetAllKeyValues(ctx context.Context) (map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var kvMap map[string]string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbname))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		kvMap = make(map[string]string, bucket.Stats().KeyN)
		return bucket.ForEach(func(k, v []byte) error {
			kvMap[string(k)] = string(v)
			return nil
		})
	})
	return kvMap, err
}
