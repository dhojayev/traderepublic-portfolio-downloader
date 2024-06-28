package document

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
)

var DefaultSet = wire.NewSet(
	NewModelBuilder,
	NewDownloader,
	NewDateResolver,
	ProvideDocumentRepository,

	wire.Bind(new(ModelBuilderInterface), new(ModelBuilder)),
	wire.Bind(new(DownloaderInterface), new(Downloader)),
	wire.Bind(new(DateResolverInterface), new(DateResolver)),
	wire.Bind(new(RepositoryInterface), new(*database.Repository[*Model])),
)

func ProvideDocumentRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Model], error) {
	return database.NewRepository[*Model](db, logger)
}
