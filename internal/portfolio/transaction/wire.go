//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package transaction

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

func ProvideHandler(
	responseReader reader.Interface,
	responseWriter writer.Interface,
	dbConnection *gorm.DB,
	logger *log.Logger,
) (Handler, error) {
	wire.Build(
		transactions.DefaultSet,
		details.DefaultSet,
		filesystem.CSVSet,
		document.DefaultSet,
		NewModelBuilderFactory,
		NewCSVEntryFactory,
		ProvideTransactionRepository,
		NewProcessor,
		NewHandler,

		wire.Bind(new(ProcessorInterface), new(Processor)),
		wire.Bind(new(ModelBuilderFactoryInterface), new(ModelBuilderFactory)),
		wire.Bind(new(RepositoryInterface), new(*database.Repository[*Model])),
	)

	return Handler{}, nil
}

func ProvideTransactionRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Model], error) {
	return database.NewRepository[*Model](db, logger)
}

func ProvideInstrumentRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Instrument], error) {
	return database.NewRepository[*Instrument](db, logger)
}
