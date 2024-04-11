package transaction

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepository(database *gorm.DB, logger *log.Logger) (*Repository, error) {
	if err := database.AutoMigrate(&Transaction{}, &Instrument{}); err != nil {
		return nil, fmt.Errorf("could not auto-migrate: %w", err)
	}

	return &Repository{
		db:     database,
		logger: logger,
	}, nil
}

func (r *Repository) Create(transaction *Transaction) error {
	r.logger.WithField("model", transaction).Trace("saving to db")

	result := r.db.Create(transaction)
	if result.Error != nil {
		return fmt.Errorf("failed creating: %w", result.Error)
	}

	r.logger.WithFields(log.Fields{
		"row affected": result.RowsAffected,
	}).Debug("saved entry to db")

	return nil
}
