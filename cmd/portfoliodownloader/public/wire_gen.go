// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/activity"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

func ProvideNonWritingApp(logger *logrus.Logger) (portfoliodownloader.App, error) {
	nilWriter := writer.NewNilWriter()
	reader, err := websocket.ProvideReader(nilWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	handler, err := activity.ProvideHandler(reader, nilWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	db, err := database.NewSQLiteInMemory(logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	transactionHandler, err := transaction.ProvideHandler(reader, nilWriter, db, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	app := portfoliodownloader.NewApp(handler, transactionHandler, logger)
	return app, nil
}

func ProvideWritingApp(logger *logrus.Logger) (portfoliodownloader.App, error) {
	jsonWriter := filesystem.NewJSONWriter(logger)
	reader, err := websocket.ProvideReader(jsonWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	handler, err := activity.ProvideHandler(reader, jsonWriter, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	db, err := database.NewSQLiteInMemory(logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	transactionHandler, err := transaction.ProvideHandler(reader, jsonWriter, db, logger)
	if err != nil {
		return portfoliodownloader.App{}, err
	}
	app := portfoliodownloader.NewApp(handler, transactionHandler, logger)
	return app, nil
}

// wire.go:

var (
	DefaultSet = wire.NewSet(websocket.ProvideReader, database.SqliteInMemorySet, activity.ProvideHandler, transaction.ProvideHandler, portfoliodownloader.NewApp, wire.Bind(new(reader.Interface), new(*websocket.Reader)), wire.Bind(new(activity.HandlerInterface), new(activity.Handler)), wire.Bind(new(transaction.HandlerInterface), new(transaction.Handler)))

	NonWritingSet = wire.NewSet(
		DefaultSet, writer.NilSet,
	)

	WritingSet = wire.NewSet(
		DefaultSet, filesystem.JSONWriterSet,
	)
)
