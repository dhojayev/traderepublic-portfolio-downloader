package timelinedetails

import (
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Handler struct {
	eventBus   *bus.EventBus
	normalizer *Normalizer
}

func NewHandler(eventBus *bus.EventBus, normalizer *Normalizer) *Handler {
	return &Handler{
		eventBus:   eventBus,
		normalizer: normalizer,
	}
}

func (h *Handler) Handle(event bus.Event) {
	var details traderepublic.TimelineDetailsJson

	err := details.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal timeline detail", "id", event.ID, "error", err)
	}

	var header traderepublic.HeaderSection

	err = details.Section(&header)
	if err != nil {
		slog.Error("failed to get header section", "err", err)
	}

	if header.Action == nil {
		return
	}

	isin := header.Action.Payload

	h.eventBus.Publish(bus.NewEvent(
		bus.TopicInstrumentFetch,
		isin,
		[]byte(isin),
	))

	err = h.normalizer.Normalize(details)
	if err != nil {
		slog.Error("error during normalization", "error", err)
	}
}
