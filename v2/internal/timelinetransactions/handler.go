package timelinetransactions

import (
	"context"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
)

type Handler struct {
	eventBus  *bus.EventBus
	msgClient message.ClientInterface
	cache     *gocache.Cache
}

func NewHandler(eventBus *bus.EventBus, msgClient message.ClientInterface, cache *gocache.Cache) *Handler {
	return &Handler{
		eventBus:  eventBus,
		msgClient: msgClient,
		cache:     cache,
	}
}

func (h *Handler) Handle(event bus.Event) {
	var transactions traderepublic.TimelineTransactionsJson

	err := transactions.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("invalid event data type", "expected", "traderepublic.TimelineTransactionsSchemaJson", "got", event.Data)
	}

	for _, transaction := range transactions.Items {
		err = h.cache.Add(string(transaction.Id), transaction, gocache.NoExpiration)
		if err != nil {
			slog.Warn("failed to add timeline transaction to cache", "id", transaction.Id)
		}

		err := h.msgClient.SubscribeToTimelineDetailV2(context.Background(), transaction.Id)
		if err != nil {
			slog.Error("failed to subscribe to timeline detail", "error", err, "transaction_id", transaction.Id)

			continue
		}
	}
}
