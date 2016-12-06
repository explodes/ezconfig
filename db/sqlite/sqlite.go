package db

import (
	"database/sql"
	"errors"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/db/registry"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteDbType = "sqlite3"
)

func init() {
	registry.Register(sqliteDbType, initDb, validateDb)
}

func validateDb(conf *ezconfig.DbConfig) error {
	if conf.Database.Host == "" {
		return errors.New("Host not specified")
	}
	return nil
}

func initDb(conf *ezconfig.DbConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", conf.Database.Host)
}
