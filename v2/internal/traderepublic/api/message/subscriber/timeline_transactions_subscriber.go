package subscriber

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type TimelineTransactionsSubscriber struct {
	name    string
	ch      <-chan []byte
	writer  writer.Writer
	counter uint
	mu      sync.Mutex
}

func NewSubscriber(name string, ch <-chan []byte, writer writer.Writer) *TimelineTransactionsSubscriber {
	return &TimelineTransactionsSubscriber{
		name:   name,
		ch:     ch,
		writer: writer,
	}
}
func (s *TimelineTransactionsSubscriber) Listen() {
	go func() {
		for data := range s.ch {
			s.mu.Lock()
			s.counter++
			num := s.counter
			s.mu.Unlock()

			slog.Info("data received", "data", string(data), "name", s.name)

			filepath := fmt.Sprintf("%s/%s", s.name, strconv.FormatUint(uint64(num), 10))

			err := s.writer.Bytes(filepath, data)
			if err == nil {
				continue
			}

			slog.Error("failed to write data", "data", string(data), "name", s.name, "error", err)
		}
	}()
}
