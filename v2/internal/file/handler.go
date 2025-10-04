package file

import (
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type RawResponseHandler struct {
	writer *writer.ResponseWriter
}

func NewRawResponseHandler(writer *writer.ResponseWriter) *RawResponseHandler {
	return &RawResponseHandler{
		writer: writer,
	}
}

func (h *RawResponseHandler) Handle(event bus.Event) {
	slog.Debug("handling event", "topic", event.Topic)

	filepath := fmt.Sprintf("%s/%s", event.Topic, event.ID)
	data := event.Data

	err := h.writer.Bytes(filepath, data.([]byte))
	if err == nil {
		return
	}

	slog.Error("failed to write data", "data", event.Data, "error", err)
}

type CSVHander struct {
	filepath string
	writer   CSVWriterInterface
}

func NewCSVHandler(filepath string, writer CSVWriterInterface) *CSVHander {
	return &CSVHander{
		filepath: filepath,
		writer:   writer,
	}
}

func (h *CSVHander) Handle(event bus.Event) {
	model, ok := event.Data.(transaction.Model)
	if !ok {
		slog.Error("invalid model received", "id", event.ID)
	}

	err := h.writer.Write(h.filepath, model)
	if err != nil {
		slog.Error("failed to write entry to csv", "id", event.ID, "filepath", h.filepath, "err", err)
	}
}
