//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package activity

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

func ProvideHandler(
	responseReader reader.Interface,
	responseWriter writer.Interface,
	logger *log.Logger,
) (Handler, error) {
	wire.Build(
		activitylog.DefaultSet,
		details.DefaultSet,
		document.DefaultSet,
		NewProcessor,
		NewHandler,

		wire.Bind(new(ProcessorInterface), new(Processor)),
	)

	return Handler{}, nil
}
