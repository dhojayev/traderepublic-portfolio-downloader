package file

import (
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type EventWriterHandler struct {
	writer *writer.ResponseWriter
}

func NewEventWriterHandler(writer *writer.ResponseWriter) *EventWriterHandler {
	return &EventWriterHandler{
		writer: writer,
	}
}

func (h *EventWriterHandler) Handle(event bus.Event) {
	slog.Debug("handling event", "topic", event.Topic)

	filepath := fmt.Sprintf("%s/%s", event.Topic, event.ID)
	data := event.Data

	err := h.writer.Bytes(filepath, data.([]byte))
	if err == nil {
		return
	}

	slog.Error("failed to write data", "data", event.Data, "error", err)
}
