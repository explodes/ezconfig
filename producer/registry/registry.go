package registry

import (
	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/producer"
)

type InitFunc func(conf *ezconfig.ProducerConfig) (producer.Producer, error)

type ValidateFunc func(conf *ezconfig.ProducerConfig) error

type ProducerFactory struct {
	Init     InitFunc
	Validate ValidateFunc
}

var registry = make(map[string]*ProducerFactory)

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

func Get(producerType string) (*ProducerFactory, bool) {
	factory, ok := registry[producerType]
	return factory, ok
}
