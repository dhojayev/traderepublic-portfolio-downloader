//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=repository.go -destination repository_mock.go -package=transaction

package transaction

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
)

type RepositoryInterface interface {
	database.RepositoryInterface[*Model]
}
