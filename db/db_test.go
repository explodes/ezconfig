package db

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
)

func TestInitProducer_sqlitedummy(t *testing.T) {
	conf := &DbConfig{
		Database: Host{
			Type: "sqlite3",
			Host: ":memory:",
		},
	}

	db, err := InitDb(conf, 0, 0)
	if err != nil {
		t.Fatalf("Error creating dummy producer: %v", err)
	}
	if db == nil {
		t.Fatal("Received nil producer")
	}
	db.Close()
}

func TestDetermineFactory(t *testing.T) {
	for _, fac := range determineFactoryCases {
		val, init, err := determineFactory(fac.databaseType)
		if err != nil && !fac.isError {
			t.Fatal("Unexpected error")
		}
		if err == nil && fac.isError {
			t.Fatal("Unexpected non-error")
		}
		sf1 := reflect.ValueOf(val)
		sf2 := reflect.ValueOf(fac.validate)
		if sf1.Pointer() != sf2.Pointer() {
			t.Fatal("Unexpected validate function")
		}
		sf1 = reflect.ValueOf(init)
		sf2 = reflect.ValueOf(fac.init)
		if sf1.Pointer() != sf2.Pointer() {
			t.Fatal("Unexpected init function")
		}
	}
}

func TestInitProducerWithRetries(t *testing.T) {
	conf := &DbConfig{}

	dummy, err := sqliteInitDb(conf)
	if err != nil {
		t.Fatal("error creating dummy")
	}
	if dummy == nil {
		t.Fatal("nil dummy")
	}

	attempts := 0
	init := func(conf *DbConfig) (*sql.DB, error) {
		attempts++
		if attempts == 3 {
			return dummy, nil
		}
		return nil, errors.New("Failed")
	}

	val, err := initDbWithRetries(&DbConfig{}, init, 10, 0)
	if err != nil {
		t.Fatalf("Error with factory: %v", err)
	}
	if val != dummy {
		t.Fatalf("Unexpected producer: %v", val)
	}
	if attempts != 3 {
		t.Fatalf("Unexpected number of attempts: %d", attempts)
	}

}
