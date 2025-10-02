package instrument

import (
	"context"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"

	"github.com/patrickmn/go-cache"
)

type Handler struct {
	msgClient message.ClientInterface
	c         *cache.Cache
}

func NewHandler(msgClient message.ClientInterface, c *cache.Cache) *Handler {
	return &Handler{
		msgClient: msgClient,
		c:         c,
	}
}

func (h *Handler) HandleFetch(event bus.Event) {
	isin := event.ID

	if isin == "" {
		slog.Error("empty isin received")

		return
	}

	_, found := h.c.Get(isin)
	if found {
		slog.Debug("instrument details found in cache", "isin", isin)

		return
	}

	err := h.msgClient.SubsribeToInstrument(context.Background(), isin)
	if err != nil {
		slog.Error("failed to subscribe to instrument", "isin", isin, "error", err)
	}
}

func (h *Handler) HandleReceived(event bus.Event) {
	isin := event.ID

	err := h.c.Add(isin, nil, cache.NoExpiration)
	if err != nil {
		slog.Error("failed to add isin to cache", "isin", isin)
	}
}
