package models

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"reflect"
)

func ErrUnsupportedAdapter(model interface{}, adapter *databaseadapters.DatabaseAdapter) error {
	name := reflect.TypeOf(model).Name()
	val := reflect.ValueOf(adapter)
	elem := val.Elem()

	var adapterName string

	if elem.IsNil() {
		adapterName = "nil"
	} else {
		adapterName = reflect.TypeOf(*adapter).Name()
	}

	return fmt.Errorf("database adapter %v is unsupported by %v", adapterName, name)
}
