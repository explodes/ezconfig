package db

import (
	"database/sql"
	"errors"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/db/registry"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// sqliteDbType is the value to use in configuration to connect to this database type
	sqliteDbType = "sqlite3"
)

// init registers the init and validation functions with the registry
func init() {
	registry.Register(sqliteDbType, initDb, validateConfig)
}

// validateConfig makes sure all the required settings are present for the database
func validateConfig(conf *ezconfig.DbConfig) error {
	if conf.Database.Host == "" {
		return errors.New("Host not specified")
	}
	return nil
}

// initDb establishes a connection with the given configuration
func initDb(conf *ezconfig.DbConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", conf.Database.Host)
}
