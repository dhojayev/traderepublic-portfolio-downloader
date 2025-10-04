package timelinedetails

import (
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

// Handler struct manages the handling of timeline details events.
type Handler struct {
	eventBus *bus.EventBus // EventBus to publish events
}

// NewHandler creates a new instance of Handler with the provided EventBus and Normalizer.
func NewHandler(eventBus *bus.EventBus) *Handler {
	return &Handler{
		eventBus: eventBus,
	}
}

// Handle processes an incoming bus event containing timeline details.
func (h *Handler) Handle(event bus.Event) {
	var details traderepublic.TimelineDetailsJson

	// Unmarshal the JSON data from the event into a TimelineDetailsJson struct
	err := details.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal timeline detail", "id", event.ID, "error", err)
		return
	}

	// Extract the header section from the timeline details
	header, err := details.SectionHeader()
	if err != nil {
		slog.Error("failed to get header section", "err", err)
		return
	}

	// Check if there is an action payload in the header
	if header.Action == nil {
		return
	}

	isin := header.Action.Payload

	// Publish a new event to fetch instrument details
	h.eventBus.Publish(bus.NewEvent(bus.TopicInstrumentFetch, isin, nil))
}
