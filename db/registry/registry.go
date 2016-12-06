package registry

import (
	"database/sql"

	"github.com/explodes/ezconfig"
)

// InitFunc is a function that takes database configuration and turns
// it into a database connection
type InitFunc func(conf *ezconfig.DbConfig) (*sql.DB, error)

// ValidateFunc is a function that checks configuration to see if it
// works for a given database type
type ValidateFunc func(conf *ezconfig.DbConfig) error

// DbFactory holds the requirements to validate and connect to a database
type DbFactory struct {
	Init     InitFunc
	Validate ValidateFunc
}

// registry holds the registered database types
var registry = make(map[string]*DbFactory)

// Register registers init and validation functions for a given database type
func Register(dbType string, init InitFunc, validate ValidateFunc) {
	if init == nil {
		panic("ezconfig: init function is nil")
	}
	if validate == nil {
		panic("ezconfig: validate function is nil")
	}
	if _, dup := registry[dbType]; dup {
		panic("ezconfig: Register called twice for type " + dbType)
	}
	registry[dbType] = &DbFactory{
		Init:     init,
		Validate: validate,
	}
}

// Get acquires the registered database type and returns its related init and validation functions
func Get(dbType string) (*DbFactory, bool) {
	factory, ok := registry[dbType]
	return factory, ok
}
