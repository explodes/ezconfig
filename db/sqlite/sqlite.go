package db

import (
	"database/sql"
	"errors"

	"github.com/explodes/ezconfig/db"
	"github.com/explodes/ezconfig/db/registry"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteDbType = "sqlite"
)

func init() {
	registry.Register(sqliteDbType, sqliteInitDb, sqliteValidateConfig)
}

func sqliteValidateConfig(conf *db.DbConfig) error {
	if conf.Database.Host == "" {
		return errors.New("Host not specified")
	}
	return nil
}

func sqliteInitDb(conf *db.DbConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", conf.Database.Host)
}
