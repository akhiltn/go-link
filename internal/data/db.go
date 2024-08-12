package data

import (
	"log"
)

type queries interface {
	Get(string) (string, error)
	Set(string, string) error
	Delete(string) error
	Close() error
}

func InitDB() (queries, error) {
	log.Println("Initializing DB")
	return initBoltDB()
}
