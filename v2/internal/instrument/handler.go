package instrument

import (
	"context"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"
)

type Handler struct {
	msgClient message.ClientInterface
}

func NewHandler(msgClient message.ClientInterface) *Handler {
	return &Handler{
		msgClient: msgClient,
	}
}

func (h *Handler) Handle(event bus.Event) {
	isin := event.ID

	if isin == "" {
		slog.Error("empty isin received")

		return
	}

	err := h.msgClient.SubsribeToInstrument(context.Background(), isin)
	if err != nil {
		slog.Error("failed to subscribe to instrument", "isin", isin, "error", err)
	}
}
