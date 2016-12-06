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

// InitProducer establishes a connection to a Producer with the given strategy
func InitProducer(conf *ezconfig.ProducerConfig, attempts int, wait backoff.Strategy) (producer.Producer, error) {
	// determine type
	factory, ok := registry.Get(conf.Settings.Type)
	if !ok {
		return nil, fmt.Errorf("Invalid producer type %s (was the producer type imported?)", conf.Settings.Type)
	}
	// validate
	if err := factory.Validate(conf); err != nil {
		return nil, err
	}
	return initProducerWithRetries(conf, factory.Init, attempts, wait)
}

// initProducerWithRetries attempts to connect to a producer a given number of times.
// If attempts is less than or equal to one, only one attempt will be made.
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
