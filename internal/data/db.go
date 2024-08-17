package data

import (
	"log"
)

type queries interface {
	Get(string) (string, error)
	Set(string, string) error
	Delete(string) error
	Close() error
	GetAllKeyValues() (map[string]string, error)
}

func GetDB() (queries, error) {
	log.Println("Getting DB")
	return initBoltDB()
}
