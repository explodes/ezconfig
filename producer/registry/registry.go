package registry

import (
	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/producer"
)

// InitFunc is a function that takes producer configuration and turns
// it into a Producer
type InitFunc func(conf *ezconfig.ProducerConfig) (producer.Producer, error)

// ValidateFunc is a function that checks configuration to see if it
// works for a given producer type
type ValidateFunc func(conf *ezconfig.ProducerConfig) error

// ProducerFactory holds the requirements to validate and connect to a producer
type ProducerFactory struct {
	Init     InitFunc
	Validate ValidateFunc
}

// registry holds the registered producer types
var registry = make(map[string]*ProducerFactory)

// Register registers init and validation functions for a given producer type
func Register(producerType string, init InitFunc, validate ValidateFunc) {
	if init == nil {
		panic("ezconfig: init function is nil")
	}
	if validate == nil {
		panic("ezconfig: validate function is nil")
	}
	if _, dup := registry[producerType]; dup {
		panic("ezconfig: Register called twice for type " + producerType)
	}
	registry[producerType] = &ProducerFactory{
		Init:     init,
		Validate: validate,
	}
}

// Get acquires the registered producer type and returns its related init and validation functions
func Get(producerType string) (*ProducerFactory, bool) {
	factory, ok := registry[producerType]
	return factory, ok
}
