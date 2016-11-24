package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/explodes/ezconfig/backoff"
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
	switch databaseType {
	case "sqlite3":
		return sqliteValidateConfig, sqliteInitDb, nil
	case "postgres":
		return postgresValidateConfig, postgresInitDb, nil
	default:
		return nil, nil, fmt.Errorf("Unsupported database type %q", databaseType)
	}
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
	return db, err
}
