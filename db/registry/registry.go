package registry

import (
	"database/sql"

	"github.com/explodes/ezconfig"
)

type InitFunc func(conf *ezconfig.DbConfig) (*sql.DB, error)

type ValidateFunc func(conf *ezconfig.DbConfig) error

type DbFactory struct {
	Init     InitFunc
	Validate ValidateFunc
}

var registry = make(map[string]*DbFactory)

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

func Get(dbType string) (*DbFactory, bool) {
	factory, ok := registry[dbType]
	return factory, ok
}
