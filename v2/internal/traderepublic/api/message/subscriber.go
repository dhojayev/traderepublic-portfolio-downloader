package message

import (
	"log/slog"
)

type SubscriberInterface interface {
	Listen()
}

type Subscriber struct {
	name string
	ch   <-chan []byte
	log  *slog.Logger
}

func NewSubscriber(name string, ch <-chan []byte, log *slog.Logger) *Subscriber {
	return &Subscriber{
		name: name,
		ch:   ch,
		log:  log,
	}
}
func (s *Subscriber) Listen() {
	go func() {
		for data := range s.ch {
			s.log.Info("data received", "data", string(data), "name", s.name)
		}
	}()
}
