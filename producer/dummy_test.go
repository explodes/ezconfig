package producer

import "testing"

func TestDummyProducer_Publish(t *testing.T) {
	dummy := dummyProducer{}
	dummy.Publish("foo", "bar")
}

func TestDummyProducer_Close(t *testing.T) {
	dummy := dummyProducer{}
	dummy.Close()
}
