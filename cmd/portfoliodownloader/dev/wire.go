// go:build wireinject
//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

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
	log "github.com/sirupsen/logrus"
)

var (
	DefaultSet = wire.NewSet(
		activitylog.DefaultSet,
		details.DefaultSet,
		transactions.DefaultSet,
		activity.DefaultSet,
		document.DefaultSet,
		transaction.DefaultSet,
		portfoliodownloader.NewApp,
		filesystem.CSVSet,
		database.SqliteOnFilesystemSet,
	)

	RemoteSet = wire.NewSet(
		DefaultSet,
		api.DefaultSet,
		auth.DefaultSet,
		console.DefaultSet,
		websocket.DefaultSet,
		filesystem.JSONWriterSet,
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
