package database

import (
	"fmt"

	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

const DBFilename = "traderepublic.db"

func NewSQLiteOnFS() (*gorm.DB, error) {
	return newSQLite(DBFilename)
}

func NewSQLiteInMemory() (*gorm.DB, error) {
	return newSQLite(":memory:?_pragma=foreign_keys(1)")
}

func newSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite db: %w", err)
	}

	return db, nil
}
