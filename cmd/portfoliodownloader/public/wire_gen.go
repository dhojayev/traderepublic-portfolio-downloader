// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func CreateNonWritingApp(logger *logrus.Logger) (portfoliodownloader.App, error) {
	client := api.NewClient(logger)
	authClient, err := auth.NewClient(client, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	authService := console.NewAuthService(authClient)
	nilWriter := writer.NewNilWriter()
	reader, err := websocket.NewReader(authService, nilWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	transactionsClient := transactions.NewClient(reader)
	detailsClient := details.NewClient(reader)
	typeResolver := details.NewTypeResolver(logger)
	modelBuilderFactory := transaction.NewModelBuilderFactory(typeResolver, logger)
	db, err := database.NewSQLiteInMemory()
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	repository, err := ProvideTransactionRepository(db, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	csvEntryFactory := transaction.NewCSVEntryFactory(logger)
	csvReader := filesystem.NewCSVReader(logger)
	csvWriter := filesystem.NewCSVWriter(logger)
	processor := transaction.NewProcessor(modelBuilderFactory, repository, csvEntryFactory, csvReader, csvWriter, logger)
	app := portfoliodownloader.NewApp(transactionsClient, detailsClient, processor, logger)
	return app, nil
}

func CreateWritingApp(logger *logrus.Logger) (portfoliodownloader.App, error) {
	client := api.NewClient(logger)
	authClient, err := auth.NewClient(client, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	authService := console.NewAuthService(authClient)
	jsonWriter := filesystem.NewJSONWriter(logger)
	reader, err := websocket.NewReader(authService, jsonWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	transactionsClient := transactions.NewClient(reader)
	detailsClient := details.NewClient(reader)
	typeResolver := details.NewTypeResolver(logger)
	modelBuilderFactory := transaction.NewModelBuilderFactory(typeResolver, logger)
	db, err := database.NewSQLiteInMemory()
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	repository, err := ProvideTransactionRepository(db, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	csvEntryFactory := transaction.NewCSVEntryFactory(logger)
	csvReader := filesystem.NewCSVReader(logger)
	csvWriter := filesystem.NewCSVWriter(logger)
	processor := transaction.NewProcessor(modelBuilderFactory, repository, csvEntryFactory, csvReader, csvWriter, logger)
	app := portfoliodownloader.NewApp(transactionsClient, detailsClient, processor, logger)
	return app, nil
}

// wire.go:

var (
	DefaultSet = wire.NewSet(api.NewClient, auth.NewClient, console.NewAuthService, websocket.NewReader, portfoliodownloader.NewApp, transactions.NewClient, details.NewClient, details.NewTypeResolver, database.NewSQLiteInMemory, transaction.NewModelBuilderFactory, transaction.NewCSVEntryFactory, filesystem.NewCSVReader, filesystem.NewCSVWriter, transaction.NewProcessor, ProvideTransactionRepository,
		ProvideInstrumentRepository, wire.Bind(new(auth.ClientInterface), new(*auth.Client)), wire.Bind(new(console.AuthServiceInterface), new(*console.AuthService)), wire.Bind(new(portfolio.ReaderInterface), new(*websocket.Reader)), wire.Bind(new(details.TypeResolverInterface), new(details.TypeResolver)), wire.Bind(new(transaction.ModelBuilderFactoryInterface), new(transaction.ModelBuilderFactory)), wire.Bind(new(transaction.RepositoryInterface), new(*database.Repository[*transaction.Model])), wire.Bind(new(transaction.InstrumentRepositoryInterface), new(*database.Repository[*transaction.Instrument])),
	)

	NonWritingSet = wire.NewSet(
		DefaultSet, writer.NewNilWriter, wire.Bind(new(writer.Interface), new(writer.NilWriter)),
	)

	WritingSet = wire.NewSet(
		DefaultSet, filesystem.NewJSONWriter, wire.Bind(new(writer.Interface), new(*filesystem.JSONWriter)),
	)
)

func ProvideTransactionRepository(db *gorm.DB, logger *logrus.Logger) (*database.Repository[*transaction.Model], error) {
	return database.NewRepository[*transaction.Model](db, logger)
}

func ProvideInstrumentRepository(db *gorm.DB, logger *logrus.Logger) (*database.Repository[*transaction.Instrument], error) {
	return database.NewRepository[*transaction.Instrument](db, logger)
}
