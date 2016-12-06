package dummy

import (
	"reflect"
	"testing"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/opener"
	"github.com/explodes/ezconfig/producer/registry"
)

func TestDetermineFactory(t *testing.T) {
	factory, ok := registry.Get(dummyProducerType)
	if !ok {
		t.Fatal("Dummy factory not registered")
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

func TestDummyProducer_Publish(t *testing.T) {
	dummy := dummyProducer{}
	dummy.Publish("foo", "bar")
}

func TestDummyProducer_Close(t *testing.T) {
	dummy := dummyProducer{}
	dummy.Close()
}

func TestInitProducer_dummy(t *testing.T) {
	conf := &ezconfig.ProducerConfig{
		Settings: ezconfig.ProducerSettings{
			Type:    "dummy",
			Retries: 5,
		},
		Hosts: []ezconfig.ProducerHost{
			{Host: "dummy", Port: 0},
		},
	}
	p, err := opener.InitProducer(conf, 0, backoff.Constant(1))
	if err != nil {
		t.Fatalf("Error creating dummy producer: %v", err)
	}
	if p == nil {
		t.Fatal("Received nil producer")
	}
	p.Close()
}
