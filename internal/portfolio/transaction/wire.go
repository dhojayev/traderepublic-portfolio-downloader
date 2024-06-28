//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package transaction

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
)

var DefaultSet = wire.NewSet(
	NewModelBuilderFactory,
	NewCSVEntryFactory,
	NewProcessor,
	NewHandler,
	ProvideTransactionRepository,
	ProvideInstrumentRepository,

	wire.Bind(new(ProcessorInterface), new(Processor)),
	wire.Bind(new(ModelBuilderFactoryInterface), new(ModelBuilderFactory)),
	wire.Bind(new(RepositoryInterface), new(*database.Repository[*Model])),
	wire.Bind(new(InstrumentRepositoryInterface), new(*database.Repository[*Instrument])),
	wire.Bind(new(HandlerInterface), new(Handler)),
)

func ProvideHandler(logger *log.Logger) (Handler, error) {
	wire.Build(
		transactions.DefaultSet,
		details.DefaultSet,
		websocket.DefaultSet,
		console.DefaultSet,
		auth.DefaultSet,
		api.DefaultSet,
		filesystem.JSONWriterSet,
		filesystem.CSVSet,
		document.DefaultSet,
		database.SqliteInMemorySet,
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
