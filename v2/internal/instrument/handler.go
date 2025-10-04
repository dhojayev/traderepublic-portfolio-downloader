package instrument

import (
	"context"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"

	gocache "github.com/patrickmn/go-cache"
)

type Handler struct {
	msgClient message.ClientInterface
	cache     *gocache.Cache
}

func NewHandler(msgClient message.ClientInterface, cache *gocache.Cache) *Handler {
	return &Handler{
		msgClient: msgClient,
		cache:     cache,
	}
}

func (h *Handler) HandleFetch(event bus.Event) {
	isin := event.ID

	if isin == "" {
		slog.Error("empty isin received")

		return
	}

	_, found := h.cache.Get(isin)
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

	var instr traderepublic.InstrumentJson

	err := instr.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal instrument", "isin", isin)
	}

	_ = h.cache.Add(isin, instr, gocache.NoExpiration)
}
