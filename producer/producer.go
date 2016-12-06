package producer

type Producer interface {
	Publish(topic string, message string)
	Close() error
}
