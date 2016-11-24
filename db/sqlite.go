package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

func sqliteValidateConfig(conf *DbConfig) error {
	if conf.Database.Host == "" {
		return errors.New("Host not specified")
	}
	return nil
}

func sqliteInitDb(conf *DbConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", conf.Database.Host)
}
