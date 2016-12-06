package opener

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/db/registry"
)

// InitDb establishes a connection to a database with the given strategy
func InitDb(conf *ezconfig.DbConfig, attempts int, wait backoff.Strategy) (*sql.DB, error) {
	// determine type
	factory, ok := registry.Get(conf.Database.Type)
	if !ok {
		return nil, fmt.Errorf("Invalid database type %s (was the database type imported?)", conf.Database.Type)
	}
	// validate
	if err := factory.Validate(conf); err != nil {
		return nil, err
	}
	return initDbWithRetries(conf, factory.Init, attempts, wait)
}

// initDbWithRetries attempts to connect to a database a given number of times.
// If attempts is less than or equal to one, only one attempt will be made.
// A "Ping" is sent to the database to test the connection.
func initDbWithRetries(conf *ezconfig.DbConfig, init registry.InitFunc, attempts int, wait backoff.Strategy) (*sql.DB, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var db *sql.DB
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		db, err = init(conf)
		if err != nil {
			log.Printf("Unable to connect to database (attempt %d of %d): %v", attempt+1, attempts, err)
			time.Sleep(wait.Duration(attempt))
			continue
		}
		if err = db.Ping(); err != nil {
			log.Printf("Unable to connect to database (attempt %d of %d) (%v)", attempt+1, attempts, err)
			time.Sleep(wait.Duration(attempt))
			continue
		}
		break
	}
	if err != nil {
		log.Printf("Unable to connect to database after %d tries", attempts)
		return nil, err
	}
	db.SetMaxOpenConns(conf.Database.MaxConnections)
	return db, nil
}
