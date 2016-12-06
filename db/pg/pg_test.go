package db

import (
	"reflect"
	"testing"

	"github.com/explodes/ezconfig/db/registry"
)

func TestDetermineFactory(t *testing.T) {
	factory, ok := registry.Get(pgDbType)
	if !ok {
		t.Fatal("Postgres factory not registered")
	}
	sf1 := reflect.ValueOf(initDb)
	sf2 := reflect.ValueOf(factory.Init)
	if sf1.Pointer() != sf2.Pointer() {
		t.Fatal("Unexpected init function")
	}
	sf1 = reflect.ValueOf(validateDb)
	sf2 = reflect.ValueOf(factory.Validate)
	if sf1.Pointer() != sf2.Pointer() {
		t.Fatal("Unexpected validate function")
	}
}
