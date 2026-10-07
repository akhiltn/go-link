package data

import "github.com/boltdb/bolt"

// BoltDB adapts BoltDB to the short URL repository interface.
type BoltDB struct {
	db *bolt.DB
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}
