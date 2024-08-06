package imdb

import (
	databaseadapters "golang-web-core/srv/database_adapters"
	"reflect"
)

type Imdb struct {
	databaseadapters.ConnectionConfig
}

func NewImdbAdapter() Imdb {
	return Imdb{
		ConnectionConfig: DefaultConfig(),
	}
}

func (db Imdb) Name() string {
	return reflect.TypeOf(db).Name()
}

func (db Imdb) Connection() databaseadapters.ConnectionConfig {
	return db.ConnectionConfig
}
