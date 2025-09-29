package subscriber

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type SubscriberInterface interface {
	Listen()
}

type Subscriber struct {
	name    string
	ch      <-chan []byte
	writer  writer.Writer
	log     *slog.Logger
	counter uint
	mu      sync.Mutex
}

func NewSubscriber(name string, ch <-chan []byte, writer writer.Writer, log *slog.Logger) *Subscriber {
	return &Subscriber{
		name:   name,
		ch:     ch,
		writer: writer,
		log:    log,
	}
}
func (s *Subscriber) Listen() {
	go func() {
		for data := range s.ch {
			s.mu.Lock()
			s.counter++
			num := s.counter
			s.mu.Unlock()

			s.log.Info("data received", "data", string(data), "name", s.name)

			filepath := fmt.Sprintf("%s/%s", s.name, strconv.FormatUint(uint64(num), 10))

			err := s.writer.Bytes(filepath, data)
			if err != nil {
				s.log.Error("failed to write data", "data", string(data), "name", s.name, "error", err)
			}
		}
	}()
}
