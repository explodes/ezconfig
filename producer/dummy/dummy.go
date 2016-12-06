package dummy

import (
	"log"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/producer"
	"github.com/explodes/ezconfig/producer/registry"
)

const (
	dummyProducerType = "dummy"
)

func init() {
	registry.Register(dummyProducerType, initProducer, validateConfig)
}

func validateConfig(conf *ezconfig.ProducerConfig) error {
	return nil
}

func initProducer(conf *ezconfig.ProducerConfig) (producer.Producer, error) {
	dummy := dummyProducer{}
	return &dummy, nil
}

type dummyProducer struct {
}

func (d dummyProducer) Publish(topic string, message string) {
	log.Printf("publish %q -> %s", topic, message)
}

func (d dummyProducer) Close() error {
	return nil
}
