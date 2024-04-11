//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package main

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
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
		transaction.NewTypeResolver,
		transaction.NewBuilder,
		transaction.NewCSVEntryFactory,
		filesystem.NewCSVReader,
		filesystem.NewCSVWriter,
		transaction.NewProcessor,

		wire.Bind(new(transaction.BuilderInterface), new(transaction.Builder)),
		wire.Bind(new(filesystem.FactoryInterface), new(transaction.CSVEntryFactory)),
	)

	RemoteSet = wire.NewSet(
		DefaultSet,
		api.NewClient,
		auth.NewClient,
		filesystem.NewJSONWriter,
		websocket.NewReader,

		wire.Bind(new(auth.ClientInterface), new(*auth.Client)),
		wire.Bind(new(writer.Interface), new(filesystem.JSONWriter)),
		wire.Bind(new(portfolio.ReaderInterface), new(*websocket.Reader)),
	)

	LocalSet = wire.NewSet(
		DefaultSet,
		writer.NewNilWriter,
		filesystem.NewJSONReader,

		wire.Bind(new(writer.Interface), new(writer.NilWriter)),
		wire.Bind(new(portfolio.ReaderInterface), new(filesystem.JSONReader)),
	)
)

func CreateLocalApp(baseDir string, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(LocalSet)

	return portfoliodownloader.App{}, nil
}

func CreateRemoteApp(phoneNumber auth.PhoneNumber, pin auth.Pin, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(RemoteSet)

	return portfoliodownloader.App{}, nil
}
