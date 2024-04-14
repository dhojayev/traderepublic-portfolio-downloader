package database

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(model any) error
}

type Repository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepository(model any, database *gorm.DB, logger *log.Logger) (*Repository, error) {
	repository := &Repository{
		db:     database,
		logger: logger,
	}

	if err := repository.db.AutoMigrate(model); err != nil {
		return nil, fmt.Errorf("could not auto-migrate: %w", err)
	}

	repository.logger.WithField("model", model).Trace("initialized repository for model")

	return repository, nil
}

func (r *Repository) Create(model any) error {
	r.logger.WithField("model", model).Trace("saving to db")

	result := r.db.Create(model)
	if result.Error != nil {
		return fmt.Errorf("failed creating: %w", result.Error)
	}

	r.logger.WithFields(log.Fields{
		"row affected": result.RowsAffected,
	}).Debug("saved entry to db")

	return nil
}
