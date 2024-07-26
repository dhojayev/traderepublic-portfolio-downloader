package document

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
)

type RepositoryInterface interface {
	database.RepositoryInterface[*Model]
}
