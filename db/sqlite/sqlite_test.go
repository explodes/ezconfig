package db

import (
	"reflect"
	"testing"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/db/registry"
	"github.com/explodes/ezconfig/opener"
)

func TestDetermineFactory(t *testing.T) {
	factory, ok := registry.Get(sqliteDbType)
	if !ok {
		t.Fatal("Sqlite factory not registered")
	}
	sf1 := reflect.ValueOf(initDb)
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

func TestInitDatbase_sqlitedummy(t *testing.T) {
	conf := &ezconfig.DbConfig{
		Database: ezconfig.DbHost{
			Type: "sqlite3",
			Host: ":memory:",
		},
	}

	db, err := opener.InitDb(conf, 0, backoff.Constant(1))
	if err != nil {
		t.Fatalf("Error creating dummy database: %v", err)
	}
	if db == nil {
		t.Fatal("Received nil database")
	}
	db.Close()
}
