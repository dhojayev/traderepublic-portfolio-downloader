//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package websocket

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var DefaultSet = wire.NewSet(
	NewReader,

	wire.Bind(new(reader.Interface), new(*Reader)),
)

func ProvideReader(responseWriter writer.Interface, logger *log.Logger) (*Reader, error) {
	wire.Build(
		auth.DefaultSet,
		api.DefaultSet,
		console.DefaultSet,
		NewReader,
	)

	return &Reader{}, nil
}
