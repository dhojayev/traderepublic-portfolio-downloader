// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/activity"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
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
	activitylogClient := activitylog.NewClient(reader, logger)
	detailsClient := details.NewClient(reader, logger)
	responseNormalizer := details.NewResponseNormalizer(logger)
	dateResolver := document.NewDateResolver(logger)
	modelBuilder := document.NewModelBuilder(dateResolver, logger)
	downloader := document.NewDownloader(logger)
	processor := activity.NewProcessor(modelBuilder, downloader, logger)
	handler := activity.NewHandler(activitylogClient, detailsClient, responseNormalizer, processor, logger)
	transactionsClient := transactions.NewClient(reader, logger)
	eventTypeResolver := transactions.NewEventTypeResolver(logger)
	typeResolver := details.NewTypeResolver(logger)
	modelBuilderFactory := transaction.NewModelBuilderFactory(typeResolver, modelBuilder, logger)
	db, err := database.NewSQLiteInMemory(logger)
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
	transactionProcessor := transaction.NewProcessor(modelBuilderFactory, repository, csvEntryFactory, csvReader, csvWriter, downloader, logger)
	transactionHandler := transaction.NewHandler(transactionsClient, detailsClient, responseNormalizer, eventTypeResolver, transactionProcessor, logger)
	app := portfoliodownloader.NewApp(handler, transactionHandler, logger)
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
	activitylogClient := activitylog.NewClient(reader, logger)
	detailsClient := details.NewClient(reader, logger)
	responseNormalizer := details.NewResponseNormalizer(logger)
	dateResolver := document.NewDateResolver(logger)
	modelBuilder := document.NewModelBuilder(dateResolver, logger)
	downloader := document.NewDownloader(logger)
	processor := activity.NewProcessor(modelBuilder, downloader, logger)
	handler := activity.NewHandler(activitylogClient, detailsClient, responseNormalizer, processor, logger)
	transactionsClient := transactions.NewClient(reader, logger)
	eventTypeResolver := transactions.NewEventTypeResolver(logger)
	typeResolver := details.NewTypeResolver(logger)
	modelBuilderFactory := transaction.NewModelBuilderFactory(typeResolver, modelBuilder, logger)
	db, err := database.NewSQLiteInMemory(logger)
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
	transactionProcessor := transaction.NewProcessor(modelBuilderFactory, repository, csvEntryFactory, csvReader, csvWriter, downloader, logger)
	transactionHandler := transaction.NewHandler(transactionsClient, detailsClient, responseNormalizer, eventTypeResolver, transactionProcessor, logger)
	app := portfoliodownloader.NewApp(handler, transactionHandler, logger)
	return app, nil
}

// wire.go:

var (
	DefaultSet = wire.NewSet(api.NewClient, auth.NewClient, console.NewAuthService, websocket.NewReader, portfoliodownloader.NewApp, transactions.NewClient, transactions.NewEventTypeResolver, details.NewClient, details.NewTypeResolver, database.NewSQLiteInMemory, transaction.NewModelBuilderFactory, document.NewModelBuilder, transaction.NewCSVEntryFactory, filesystem.NewCSVReader, filesystem.NewCSVWriter, transaction.NewProcessor, document.NewDownloader, document.NewDateResolver, ProvideTransactionRepository,
		ProvideInstrumentRepository,
		ProvideDocumentRepository, activitylog.NewClient, activity.NewProcessor, activity.NewHandler, transaction.NewHandler, details.NewResponseNormalizer, wire.Bind(new(auth.ClientInterface), new(*auth.Client)), wire.Bind(new(console.AuthServiceInterface), new(*console.AuthService)), wire.Bind(new(reader.Interface), new(*websocket.Reader)), wire.Bind(new(transactions.ClientInterface), new(transactions.Client)), wire.Bind(new(transactions.EventTypeResolverInterface), new(transactions.EventTypeResolver)), wire.Bind(new(details.ClientInterface), new(details.Client)), wire.Bind(new(details.TypeResolverInterface), new(details.TypeResolver)), wire.Bind(new(transaction.ProcessorInterface), new(transaction.Processor)), wire.Bind(new(transaction.ModelBuilderFactoryInterface), new(transaction.ModelBuilderFactory)), wire.Bind(new(document.ModelBuilderInterface), new(document.ModelBuilder)), wire.Bind(new(transaction.RepositoryInterface), new(*database.Repository[*transaction.Model])), wire.Bind(new(transaction.InstrumentRepositoryInterface), new(*database.Repository[*transaction.Instrument])), wire.Bind(new(document.DownloaderInterface), new(document.Downloader)), wire.Bind(new(document.DateResolverInterface), new(document.DateResolver)), wire.Bind(new(document.RepositoryInterface), new(*database.Repository[*document.Model])), wire.Bind(new(filesystem.CSVReaderInterface), new(filesystem.CSVReader)), wire.Bind(new(filesystem.CSVWriterInterface), new(filesystem.CSVWriter)), wire.Bind(new(activitylog.ClientInterface), new(activitylog.Client)), wire.Bind(new(activity.ProcessorInterface), new(activity.Processor)), wire.Bind(new(activity.HandlerInterface), new(activity.Handler)), wire.Bind(new(transaction.HandlerInterface), new(transaction.Handler)), wire.Bind(new(details.ResponseNormalizerInterface), new(details.ResponseNormalizer)),
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

func ProvideDocumentRepository(db *gorm.DB, logger *logrus.Logger) (*database.Repository[*document.Model], error) {
	return database.NewRepository[*document.Model](db, logger)
}
