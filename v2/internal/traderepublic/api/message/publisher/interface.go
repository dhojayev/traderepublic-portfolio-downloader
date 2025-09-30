package publisher

type Interface interface {
	Subscribe(topic string) <-chan []byte
	Publish(msg []byte, topic string)
	Close(topic string)
}
