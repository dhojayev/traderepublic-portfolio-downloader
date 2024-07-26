//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

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
	log "github.com/sirupsen/logrus"
)

var (
	DefaultSet = wire.NewSet(
		database.SqliteOnFilesystemSet,
		activity.ProvideHandler,
		transaction.ProvideHandler,
		portfoliodownloader.NewApp,

		wire.Bind(new(activity.HandlerInterface), new(activity.Handler)),
		wire.Bind(new(transaction.HandlerInterface), new(transaction.Handler)),
	)

	RemoteSet = wire.NewSet(
		DefaultSet,
		websocket.ProvideReader,
		filesystem.JSONWriterSet,

		wire.Bind(new(reader.Interface), new(*websocket.Reader)),
	)

	LocalSet = wire.NewSet(
		DefaultSet,
		reader.DefaultSet,
		writer.NilSet,
	)
)

func ProvideLocalApp(baseDir string, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(LocalSet)

	return portfoliodownloader.App{}, nil
}

func ProvideRemoteApp(logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(RemoteSet)

	return portfoliodownloader.App{}, nil
}
