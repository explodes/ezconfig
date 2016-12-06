package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/db/registry"
	_ "github.com/lib/pq"
)

const (
	// pgDbType is the value to use in configuration to connect to this database type
	pgDbType = "postgres"
)

// init registers the init and validation functions with the registry
func init() {
	registry.Register(pgDbType, initDb, validateDb)
}

// getConnectionString builds a connection string from the supplied configuration
func getConnectionString(conf *ezconfig.DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.DbName, conf.Database.Ssl)
}

// validateDb makes sure all the required settings are present for the database
func validateDb(conf *ezconfig.DbConfig) error {
	if conf.Database.Host == "" {
		return errors.New("Host not specified")
	}
	if conf.Database.User == "" {
		return errors.New("User not specified")
	}
	if conf.Database.Port == 0 {
		return errors.New("Port not specified")
	}
	if conf.Database.DbName == "" {
		return errors.New("Database not specified")
	}
	if conf.Database.Ssl == "" {
		return errors.New("Ssl not specified")
	}
	if conf.Database.Password == "" {
		return errors.New("Password not specified")
	}
	return nil
}

// initDb establishes a connection with the given configuration
func initDb(conf *ezconfig.DbConfig) (*sql.DB, error) {
	connStr := getConnectionString(conf)
	return sql.Open("postgres", connStr)
}
