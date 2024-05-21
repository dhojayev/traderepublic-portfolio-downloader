//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package main

import (
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"

	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultSet = wire.NewSet(
		portfoliodownloader.NewApp,
		transactions.NewClient,
		details.NewClient,
		details.NewTypeResolver,
		transaction.NewModelBuilderFactory,
		database.NewSQLiteOnFS,
		transaction.NewCSVEntryFactory,
		filesystem.NewCSVReader,
		filesystem.NewCSVWriter,
		transaction.NewProcessor,
		ProvideTransactionRepository,
		ProvideInstrumentRepository,

		wire.Bind(new(details.TypeResolverInterface), new(details.TypeResolver)),
		wire.Bind(new(transaction.ModelBuilderFactoryInterface), new(transaction.ModelBuilderFactory)),
		wire.Bind(new(transaction.RepositoryInterface), new(*database.Repository[*transaction.Model])),
		wire.Bind(new(transaction.InstrumentRepositoryInterface), new(*database.Repository[*transaction.Instrument])),
	)

	RemoteSet = wire.NewSet(
		DefaultSet,
		api.NewClient,
		auth.NewClient,
		console.NewAuthService,
		websocket.NewReader,
		filesystem.NewJSONWriter,

		wire.Bind(new(auth.ClientInterface), new(*auth.Client)),
		wire.Bind(new(console.AuthServiceInterface), new(*console.AuthService)),
		wire.Bind(new(portfolio.ReaderInterface), new(*websocket.Reader)),

		wire.Bind(new(writer.Interface), new(*filesystem.JSONWriter)),
	)

	LocalSet = wire.NewSet(
		DefaultSet,
		writer.NewNilWriter,
		filesystem.NewJSONReader,

		wire.Bind(new(writer.Interface), new(writer.NilWriter)),
		wire.Bind(new(portfolio.ReaderInterface), new(*filesystem.JSONReader)),
	)
)

func CreateLocalApp(baseDir string, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(LocalSet)

	return portfoliodownloader.App{}, nil
}

func CreateRemoteApp(logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(RemoteSet)

	return portfoliodownloader.App{}, nil
}

func ProvideTransactionRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*transaction.Model], error) {
	return database.NewRepository[*transaction.Model](db, logger)
}

func ProvideInstrumentRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*transaction.Instrument], error) {
	return database.NewRepository[*transaction.Instrument](db, logger)
}
