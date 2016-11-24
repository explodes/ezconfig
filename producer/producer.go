package producer

import (
	"fmt"
	"log"
	"time"

	"github.com/explodes/ezconfig/backoff"
)

type Producer interface {
	Publish(topic string, message string)
	Close()
}

type initProducerFunc func(conf *ProducerConfig) (Producer, error)
type validateConfigFunc func(conf *ProducerConfig) error

func InitProducer(conf *ProducerConfig, attempts int, wait backoff.Strategy) (Producer, error) {

	// determine type
	validate, init, err := determineFactory(conf.Settings.Type)
	if err != nil {
		return nil, err
	}

	// validate
	if err := validate(conf); err != nil {
		return nil, err
	}

	return initProducerWithRetries(conf, init, attempts, wait)
}

func determineFactory(producerType string) (validateConfigFunc, initProducerFunc, error) {
	switch producerType {
	case "dummy":
		return dummyValidateConfig, dummyInitProducer, nil
	case "kafka":
		return kafkaValidateConfig, kafkaInitProducer, nil
	default:
		return nil, nil, fmt.Errorf("Unsupported producer type %q", producerType)
	}
}

func initProducerWithRetries(conf *ProducerConfig, init initProducerFunc, attempts int, wait backoff.Strategy) (Producer, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var producer Producer
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		producer, err = init(conf)
		if err != nil {
			log.Printf("Unable to create producer (attempt %d of %d)", attempt+1, attempts)
			time.Sleep(wait.Duration(attempt))
			continue
		}
		break
	}
	if err != nil {
		log.Printf("Unable to create producer after %d tries", attempts)
		return nil, err
	}
	return producer, err
}
