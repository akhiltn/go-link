package data

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boltdb/bolt"
)

type boltDB struct {
	db *bolt.DB
}

var (
	instance *boltDB
	appPath  = ""
)

const (
	dbname = "DB"
)

func newBoltDB(dbPath string) (*boltDB, error) {
	db, err := bolt.Open(os.ExpandEnv(dbPath), 0600, nil)
	if err != nil {
		log.Printf("Error opening BoltDB: %v", err)
		return nil, err
	}
	return &boltDB{db: db}, nil
}

func (b *boltDB) Close() error {
	return b.db.Close()
}

func initBoltDB() (*boltDB, error) {
	if instance != nil {
		return instance, nil
	}
	var err error
	log.Println("Initializing BoltDB")
	var dbPath string
	if runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Println("Error getting home directory:", err)
			return nil, err
		}
		dbPath = filepath.Join(homeDir, ".quick-url", "quick-url.db")
		appPath = filepath.Dir(dbPath)
		if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
		}
	} else {
		dbPath = "quick-url.db"
	}
	instance, err = newBoltDB(dbPath)
	if err != nil {
		log.Printf("Failed to initialize BoltDB: %v", err)
	} else {
		// Ensure the bucket is created
		err = instance.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dbname))
			return err
		})
		if err != nil {
			log.Printf("Failed to create bucket: %v", err)
		}
	}
	if instance == nil {
		return nil, bolt.ErrDatabaseNotOpen
	}
	return instance, err
}

func (b *boltDB) Get(key string) (string, error) {
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

func (b *boltDB) Set(key string, value string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(dbname))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (b *boltDB) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbname))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func (b *boltDB) GetAllKeyValues() (map[string]string, error) {
	var kvMap = make(map[string]string)
	err := b.db.View(func(tx *bolt.Tx) error {
		// Iterate over each bucket
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			// Iterate over each key-value pair in the bucket
			return b.ForEach(func(k, v []byte) error {
				kvMap[string(k)] = string(v)
				return nil
			})
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	return kvMap, err
}
