package models

import databaseadapters "golang-web-core/srv/database_adapters"

type Model interface {
	Adapter() *databaseadapters.DatabaseAdapter
	Name() string
	PrimaryKey() string
	Create(object any) (any, error)
	Find(key any) (any, error)
	Where(query map[string]any) (any, error)
	All() (any, error)
	Update(key, object any) error
	UpdateWhere(query map[string]any, object any) error
	Delete(key any) error
	DeleteWhere(query map[string]any) error
}
