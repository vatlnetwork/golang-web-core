package imdb

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/util"
	"reflect"
)

type Imdb struct {
	databaseadapters.ConnectionConfig
	Data map[string][]any
}

func NewImdbAdapter(connectionConfig databaseadapters.ConnectionConfig) *Imdb {
	return &Imdb{
		ConnectionConfig: connectionConfig,
		Data:             map[string][]any{},
	}
}

func (db Imdb) Name() string {
	return reflect.TypeOf(db).Name()
}

func (db Imdb) Connection() databaseadapters.ConnectionConfig {
	return db.ConnectionConfig
}

func (db Imdb) TestConnection() error {
	return nil
}

func (db *Imdb) ApplyConfig(config databaseadapters.ConnectionConfig) {
	db.ConnectionConfig = config
}

func (db *Imdb) Insert(modelName string, object any) {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	collection = append(collection, object)
	db.Data[modelName] = collection
}

func (db *Imdb) GetAll(modelName string) []any {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	return collection
}

func (db *Imdb) Find(modelName, key string, value any) (any, error) {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	for _, item := range collection {
		json := util.StructToMap(item)
		if json[key] == value {
			return item, nil
		}
	}

	return nil, fmt.Errorf("unable to find a %v with %v = %v", modelName, key, value)
}

func (db *Imdb) Query(modelName string, query map[string]any) []any {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	results := []any{}
	for _, item := range collection {
		json := util.StructToMap(item)
		matches := true
		for key := range query {
			if json[key] != query[key] {
				matches = false
			}
		}
		if matches {
			results = append(results, item)
		}
	}
	return results
}

func (db *Imdb) Update(modelName, primaryKey string, keyValue any, object any) error {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	updatedCollection := []any{}
	found := false
	for _, item := range collection {
		json := util.StructToMap(item)
		if json[primaryKey] == keyValue {
			updatedCollection = append(updatedCollection, object)
			found = true
		} else {
			updatedCollection = append(updatedCollection, item)
		}
	}
	if !found {
		return fmt.Errorf("unable to find a %v with %v = %v", modelName, primaryKey, keyValue)
	}
	db.Data[modelName] = updatedCollection
	return nil
}

func (db *Imdb) Delete(modelName, primaryKey string, keyValue any) error {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []any{}
	}
	updatedCollection := []any{}
	found := false
	for _, item := range collection {
		json := util.StructToMap(item)
		if json[primaryKey] == keyValue {
			found = true
		} else {
			updatedCollection = append(updatedCollection, item)
		}
	}
	if !found {
		return fmt.Errorf("unable to find a %v with %v = %v", modelName, primaryKey, keyValue)
	}
	db.Data[modelName] = updatedCollection
	return nil
}
