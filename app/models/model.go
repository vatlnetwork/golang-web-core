package models

import databaseadapters "golang-web-core/srv/database_adapters"

type Model interface {
	Adapter() *databaseadapters.DatabaseAdapter
	Name() string
	PrimaryKey() string
	Create(object interface{}) (interface{}, error)
	Find(key interface{}) (interface{}, error)
	Where(query map[string]interface{}) (interface{}, error)
	All() (interface{}, error)
	Update(key, object interface{}) error
	UpdateWhere(query map[string]interface{}, object interface{}) error
	Delete(key interface{}) error
	DeleteWhere(query map[string]interface{}) error
}
