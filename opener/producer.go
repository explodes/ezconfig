package opener

import (
	"fmt"
	"log"
	"time"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/producer"
	"github.com/explodes/ezconfig/producer/registry"
)

func InitProducer(conf *ezconfig.ProducerConfig, attempts int, wait backoff.Strategy) (producer.Producer, error) {

	// determine type
	validate, init, err := determineProducerFactory(conf.Settings.Type)
	if err != nil {
		return nil, err
	}

	// validate
	if err := validate(conf); err != nil {
		return nil, err
	}

	return initProducerWithRetries(conf, init, attempts, wait)
}

func determineProducerFactory(producerType string) (registry.ValidateFunc, registry.InitFunc, error) {
	factory, ok := registry.Get(producerType)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid producer type %s (was the database type imported?)", producerType)
	}
	return factory.Validate, factory.Init, nil
}

func initProducerWithRetries(conf *ezconfig.ProducerConfig, init registry.InitFunc, attempts int, wait backoff.Strategy) (producer.Producer, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var p producer.Producer
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		p, err = init(conf)
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
	return p, err
}
