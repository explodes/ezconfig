package producer

import (
	"log"
)

func dummyValidateConfig(conf *ProducerConfig) error {
	return nil
}

func dummyInitProducer(conf *ProducerConfig) (Producer, error) {
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
