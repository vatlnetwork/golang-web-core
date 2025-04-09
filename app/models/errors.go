package models

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/srverr"
	"reflect"
)

func ErrUnsupportedAdapter(model any, adapter *databaseadapters.DatabaseAdapter) srverr.ServerError {
	name := reflect.TypeOf(model).Name()
	val := reflect.ValueOf(adapter)
	elem := val.Elem()

	var adapterName string

	if elem.IsNil() {
		adapterName = "nil"
	} else {
		adapterName = reflect.TypeOf(*adapter).Name()
	}

	return srverr.New(fmt.Sprintf("database adapter %v is unsupported by %v", adapterName, name), 500)
}
