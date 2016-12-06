package kafka

import (
	"reflect"
	"testing"

	"github.com/explodes/ezconfig/producer/registry"
)

func TestDetermineFactory(t *testing.T) {
	factory, ok := registry.Get(kafkaProducerType)
	if !ok {
		t.Fatal("Kafka factory not registered")
	}
	sf1 := reflect.ValueOf(initProducer)
	sf2 := reflect.ValueOf(factory.Init)
	if sf1.Pointer() != sf2.Pointer() {
		t.Fatal("Unexpected init function")
	}
	sf1 = reflect.ValueOf(validateConfig)
	sf2 = reflect.ValueOf(factory.Validate)
	if sf1.Pointer() != sf2.Pointer() {
		t.Fatal("Unexpected validate function")
	}
}
