//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package main

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/activity"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var (
	DefaultSet = wire.NewSet(
		websocket.ProvideReader,
		database.SqliteInMemorySet,
		activity.ProvideHandler,
		transaction.ProvideHandler,
		portfoliodownloader.NewApp,

		wire.Bind(new(reader.Interface), new(*websocket.Reader)),
		wire.Bind(new(activity.HandlerInterface), new(activity.Handler)),
		wire.Bind(new(transaction.HandlerInterface), new(transaction.Handler)),
	)

	NonWritingSet = wire.NewSet(
		DefaultSet,
		writer.NilSet,
	)

	WritingSet = wire.NewSet(
		DefaultSet,
		filesystem.JSONWriterSet,
	)
)

func ProvideNonWritingApp(logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(NonWritingSet)

	return portfoliodownloader.App{}, nil
}

func ProvideWritingApp(logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(WritingSet)

	return portfoliodownloader.App{}, nil
}
