package database

import (
	"fmt"

	"gorm.io/driver/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

const DBFilename = "traderepublic.db"

func NewSQLiteOnFS() (*gorm.DB, error) {
	return newSQLite(DBFilename)
}

func NewSQLiteInMemory() (*gorm.DB, error) {
	return newSQLite("file::memory:?cache=shared")
}

func newSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite db: %w", err)
	}

	return db, nil
}
