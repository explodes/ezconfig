package registry

type DbFactory struct {
	Init     interface{}
	Validate interface{}
}

var registry = make(map[string]*DbFactory)

func Register(dbType string, init interface{}, validate interface{}) {
	registry[dbType] = &DbFactory{
		Init:     init,
		Validate: validate,
	}
}

func Get(dbType string) (*DbFactory, bool) {
	factory, ok := registry[dbType]
	return factory, ok
}