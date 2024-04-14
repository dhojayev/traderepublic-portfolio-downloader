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
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var (
	DefaultSet = wire.NewSet(
		portfoliodownloader.NewApp,
		transactions.NewClient,
		details.NewClient,
		transaction.NewTypeResolver,
		database.NewSQLiteOnFS,
		transaction.NewRepository,
		transaction.NewBuilder,
		transaction.NewCSVEntryFactory,
		filesystem.NewCSVReader,
		filesystem.NewCSVWriter,
		transaction.NewProcessor,
		api.NewClient,
		auth.NewClient,
		websocket.NewReader,

		wire.Bind(new(auth.ClientInterface), new(*auth.Client)),
		wire.Bind(new(portfolio.ReaderInterface), new(*websocket.Reader)),
		wire.Bind(new(transaction.BuilderInterface), new(transaction.Builder)),
		wire.Bind(new(transaction.RepositoryInterface), new(*transaction.Repository)),
	)

	NonWritingSet = wire.NewSet(
		DefaultSet,
		writer.NewNilWriter,

		wire.Bind(new(writer.Interface), new(writer.NilWriter)),
	)

	WritingSet = wire.NewSet(
		DefaultSet,
		filesystem.NewJSONWriter,

		wire.Bind(new(writer.Interface), new(filesystem.JSONWriter)),
	)
)

func CreateNonWritingApp(phoneNumber auth.PhoneNumber, pin auth.Pin, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(NonWritingSet)

	return portfoliodownloader.App{}, nil
}

func CreateWritingApp(phoneNumber auth.PhoneNumber, pin auth.Pin, logger *log.Logger) (portfoliodownloader.App, error) {
	wire.Build(WritingSet)

	return portfoliodownloader.App{}, nil
}
