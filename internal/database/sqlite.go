package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DBFilename = "traderepublic.db"

func NewSQLiteOnFS(logger *log.Logger) (*gorm.DB, error) {
	return newSQLite(DBFilename, logger)
}

func NewSQLiteInMemory(logger *log.Logger) (*gorm.DB, error) {
	return newSQLite(":memory:?_pragma=foreign_keys(1)", logger)
}

func newSQLite(dsn string, logrus *log.Logger) (*gorm.DB, error) {
	logLevel := logger.Error

	if logrus.Level == log.TraceLevel {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		logrus,
		logger.Config{
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
		},
	)

	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite db: %w", err)
	}

	return database, nil
}
