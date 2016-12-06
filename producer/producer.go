package producer

// Producer is capable of publishing messages to the service which backs it
type Producer interface {

	// Publish will publish a message to a given topic
	Publish(topic string, message string)

	// Close will terminate the connection to
	// the service backing this producer
	Close() error
}
