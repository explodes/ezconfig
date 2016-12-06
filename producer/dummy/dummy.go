package dummy

import (
	"log"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/producer"
	"github.com/explodes/ezconfig/producer/registry"
)

const (
	// dummyProducerType is the value to use in configuration to connect to this producer type
	dummyProducerType = "dummy"
)

// init registers the init and validation functions with the registry
func init() {
	registry.Register(dummyProducerType, initProducer, validateConfig)
}

// validateConfig makes sure all the required settings are present for the database
func validateConfig(conf *ezconfig.ProducerConfig) error {
	return nil
}

// initProducer establishes a connection with the given configuration
func initProducer(conf *ezconfig.ProducerConfig) (producer.Producer, error) {
	dummy := dummyProducer{}
	return &dummy, nil
}

// dummyProducer is a stand-in producer that emits messages only to stdout
type dummyProducer struct {
}

// Publish "publishes" a message to stdout
func (d dummyProducer) Publish(topic string, message string) {
	log.Printf("publish %q -> %s", topic, message)
}

// Close is a no-op for dummyProducers
func (d dummyProducer) Close() error {
	return nil
}
