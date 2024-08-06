package mongo

import (
	databaseadapters "golang-web-core/srv/database_adapters"
	"reflect"
)

type Mongo struct {
	databaseadapters.ConnectionConfig
}

func NewMongoAdapter() Mongo {
	return Mongo{
		ConnectionConfig: DefaultConfig(),
	}
}

func (m Mongo) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m Mongo) Connection() databaseadapters.ConnectionConfig {
	return m.ConnectionConfig
}
