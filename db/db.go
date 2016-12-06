package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/db/registry"
)

type initDbFunc func(conf *DbConfig) (*sql.DB, error)
type validateConfigFunc func(conf *DbConfig) error

func InitDb(conf *DbConfig, attempts int, wait backoff.Strategy) (*sql.DB, error) {

	// determine type
	validate, init, err := determineFactory(conf.Database.Type)
	if err != nil {
		return nil, err
	}

	// validate
	if err := validate(conf); err != nil {
		return nil, err
	}

	return initDbWithRetries(conf, init, attempts, wait)
}

func determineFactory(databaseType string) (validateConfigFunc, initDbFunc, error) {
	factory, ok := registry.Get(databaseType)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid database type %s (was the database type imported?)", databaseType)
	}
	return factory.Validate.(validateConfigFunc), factory.Init.(initDbFunc), nil
}

func initDbWithRetries(conf *DbConfig, init initDbFunc, attempts int, wait backoff.Strategy) (*sql.DB, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var db *sql.DB
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		db, err = init(conf)
		if err != nil {
			log.Printf("Unable to connect to database (attempt %d of %d): %v", attempt + 1, attempts, err)
			time.Sleep(wait.Duration(attempt))
			continue
		}
		if err = db.Ping(); err != nil {
			log.Printf("Unable to connect to database (attempt %d of %d) (%v)", attempt + 1, attempts, err)
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
