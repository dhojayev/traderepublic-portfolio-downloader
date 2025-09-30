package timelinedetails

import (
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(event bus.Event) {
	var details traderepublic.TimelineDetailJson

	err := details.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal timeline detail", "id", event.ID, "error", err)
	}
}
