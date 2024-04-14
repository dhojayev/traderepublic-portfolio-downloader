package transaction

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	database.RepositoryInterface
}

type Repository struct {
	*database.Repository
}

func NewRepository(db *gorm.DB, logger *log.Logger) (*Repository, error) {
	baseRepo, err := database.NewRepository(&Model{}, db, logger)
	if err != nil {
		return &Repository{}, fmt.Errorf("could not create base repo: %w", err)
	}

	return &Repository{baseRepo}, nil
}
