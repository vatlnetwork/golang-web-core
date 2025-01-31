package models

import (
	databaseadapters "golang-web-core/srv/database_adapters"
	"reflect"
)

type TestModel struct {
	Adapter *databaseadapters.DatabaseAdapter
}

func NewTestModel(adapter *databaseadapters.DatabaseAdapter) TestModel {
	return TestModel{
		Adapter: adapter,
	}
}

func (m TestModel) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m TestModel) PrimaryKey() string {
	return "id"
}
