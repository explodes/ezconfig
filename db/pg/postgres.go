package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/explodes/ezconfig/db"
	_ "github.com/lib/pq"
	"github.com/explodes/ezconfig/db/registry"
)

const (
	pgDbType = "postgres"
)

func init() {
	registry.Register(pgDbType, postgresInitDb, postgresValidateConfig)
}

func getPostgresConnectionsString(conf *db.DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.DbName, conf.Database.Ssl)
}

func postgresValidateConfig(conf *db.DbConfig) error {
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

func postgresInitDb(conf *db.DbConfig) (*sql.DB, error) {
	connStr := getPostgresConnectionsString(conf)
	return sql.Open("postgres", connStr)
}
