package models

import (
	"fmt"
	"reflect"
)

// var ErrUnknownAdapter error = fmt.Errorf("the database adapter specified for TestModel is unrecognized")

func ErrUnsupportedAdapter(model interface{}, adapter interface{}) error {
	name := reflect.TypeOf(model).Name()
	adapterName := reflect.TypeOf(adapter).Name()
	return fmt.Errorf("database adapter %v is unsupported by %v", adapterName, name)
}
