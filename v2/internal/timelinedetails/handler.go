package timelinedetails

import (
	"encoding/json"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Handler struct {
	eventBus *bus.EventBus
}

func NewHandler(eventBus *bus.EventBus) *Handler {
	return &Handler{
		eventBus: eventBus,
	}
}

func (h *Handler) Handle(event bus.Event) {
	var details traderepublic.TimelineDetailsJson

	err := details.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal timeline detail", "id", event.ID, "error", err)
	}

	if details.Id == "0af6e075-8034-4ec9-9956-a6c63fe6103b" {
		slog.Info("here")
	}

	var isin string

	for _, section := range details.Sections {
		var header traderepublic.HeaderSection

		data, err := json.Marshal(section)
		if err != nil {
			continue
		}

		err = header.UnmarshalJSON(data)
		if err != nil {
			continue
		}

		if header.Action == nil {
			return
		}

		isin = header.Action.Payload

		if isin == "" {
			return
		}

		h.eventBus.Publish(bus.NewEvent(
			bus.TopicInstrumentRequested,
			isin,
			[]byte(isin),
		))

		return
	}
}
