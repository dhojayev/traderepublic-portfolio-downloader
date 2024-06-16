package database

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface[M comparable] interface {
	Create(model M) error
}

type Repository[M comparable] struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepository[M comparable](database *gorm.DB, logger *log.Logger) (*Repository[M], error) {
	var model M

	repository := &Repository[M]{
		db:     database,
		logger: logger,
	}

	if err := repository.db.AutoMigrate(model); err != nil {
		return nil, fmt.Errorf("could not auto-migrate: %w", err)
	}

	repository.logger.WithField("model", reflect.TypeOf(model)).Trace("initialized repository for model")

	return repository, nil
}

func (r *Repository[M]) Create(model M) error {
	r.logger.WithField("model", model).Trace("saving to db")

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed saving: %w", result.Error)
	}

	r.logger.WithFields(log.Fields{
		"row affected": result.RowsAffected,
	}).Debug("saved entry to db")

	return nil
}
