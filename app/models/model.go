package models

type Model interface {
	Name() string
	PrimaryKey() string
	Create(object interface{}) (interface{}, error)
	Find(key any) (interface{}, error)
	Where(query interface{}) ([]interface{}, error)
	Update(key any, object interface{}) error
	UpdateWhere(query interface{}, object interface{}) error
	Delete(key any) error
	DeleteWhere(query interface{}) error
}
