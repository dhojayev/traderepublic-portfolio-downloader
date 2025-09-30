package subscriber

import (
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
)

type MessageSubscriber struct {
	target   string
	ch       <-chan []byte
	eventBus *bus.EventBus
}

func NewMessageSubscriber(target string, ch <-chan []byte, eventBus *bus.EventBus) *MessageSubscriber {
	return &MessageSubscriber{
		target:   target,
		ch:       ch,
		eventBus: eventBus,
	}
}

func (s *MessageSubscriber) Listen() {
	go func() {
		for data := range s.ch {
			slog.Info("data received", "data", string(data), "target", s.target)

			s.eventBus.Publish(s.target, data)
		}
	}()
}
