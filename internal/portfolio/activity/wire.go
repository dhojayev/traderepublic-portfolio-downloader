//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package activity

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var DefaultSet = wire.NewSet(
	NewProcessor,
	NewHandler,

	wire.Bind(new(ProcessorInterface), new(Processor)),
	wire.Bind(new(HandlerInterface), new(Handler)),
)

func ProvideHandler(
	responseWriter writer.Interface,
	logger *log.Logger,
) (Handler, error) {
	wire.Build(
		activitylog.DefaultSet,
		details.DefaultSet,
		websocket.DefaultSet,
		console.DefaultSet,
		auth.DefaultSet,
		api.DefaultSet,
		document.DefaultSet,
		NewProcessor,
		NewHandler,

		wire.Bind(new(ProcessorInterface), new(Processor)),
	)

	return Handler{}, nil
}
