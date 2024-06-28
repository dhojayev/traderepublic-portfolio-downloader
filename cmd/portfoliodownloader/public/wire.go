//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package main

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"

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
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var (
	DefaultSet = wire.NewSet(
		activitylog.DefaultSet,
		details.DefaultSet,
		transactions.DefaultSet,
		api.DefaultSet,
		activity.DefaultSet,
		document.DefaultSet,
		transaction.DefaultSet,
		auth.DefaultSet,
		console.DefaultSet,
		websocket.DefaultSet,
		filesystem.CSVSet,
		database.SqliteInMemorySet,

		portfoliodownloader.NewApp,
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
