package subscriber

import (
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type TimelineDetailSubscriber struct {
	name   string
	ch     <-chan []byte
	writer writer.Writer
	log    *slog.Logger
}

func NewTimelineDetailSubscriber(name string, ch <-chan []byte, writer writer.Writer, log *slog.Logger) *TimelineDetailSubscriber {
	return &TimelineDetailSubscriber{
		name:   name,
		ch:     ch,
		writer: writer,
		log:    log,
	}
}

func (s *TimelineDetailSubscriber) Listen() {
	go func() {
		for data := range s.ch {
			s.log.Info("data received", "data", string(data), "name", s.name)

			err := s.writer.Bytes(s.name, data)
			if err != nil {
				s.log.Error("failed to write data", "data", string(data), "name", s.name, "error", err)
			}
		}
	}()
}
