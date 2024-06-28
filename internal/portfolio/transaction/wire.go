package transaction

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
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

func ProvideTransactionRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Model], error) {
	return database.NewRepository[*Model](db, logger)
}

func ProvideInstrumentRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Instrument], error) {
	return database.NewRepository[*Instrument](db, logger)
}
