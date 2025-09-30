package file

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type RawResponseHandler struct {
	writer  *writer.ResponseWriter
	counter uint
	mu      sync.Mutex
}

func NewRawResponseHandler(writer *writer.ResponseWriter) *RawResponseHandler {
	return &RawResponseHandler{
		writer: writer,
	}
}

func (h *RawResponseHandler) Handle(event bus.Event) {
	slog.Debug("handling event", "name", event.Name)

	h.mu.Lock()
	h.counter++
	num := h.counter
	h.mu.Unlock()

	filepath := fmt.Sprintf("%s/%s", event.Name, strconv.FormatUint(uint64(num), 10))
	data := event.Data

	err := h.writer.Bytes(filepath, data.([]byte))
	if err == nil {
		return
	}

	slog.Error("failed to write data", "data", event.Data, "name", event.Name, "error", err)
}
