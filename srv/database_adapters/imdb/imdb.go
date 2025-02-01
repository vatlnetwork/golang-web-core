package imdb

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/util"
	"reflect"
)

type Imdb struct {
	databaseadapters.ConnectionConfig
	Data map[string][]interface{}
}

func NewImdbAdapter() Imdb {
	return Imdb{
		ConnectionConfig: DefaultConfig(),
		Data:             map[string][]interface{}{},
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

func (db *Imdb) Insert(modelName string, object interface{}) {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	collection = append(collection, object)
	db.Data[modelName] = collection
}

func (db *Imdb) GetAll(modelName string) []interface{} {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	return collection
}

func (db *Imdb) Find(modelName, key string, value interface{}) (interface{}, error) {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	for _, item := range collection {
		json := util.StructToMap(item)
		if json[key] == value {
			return item, nil
		}
	}

	return nil, fmt.Errorf("Unable to find a %v with %v = %v", modelName, key, value)
}

func (db *Imdb) Query(modelName string, query map[string]interface{}) []interface{} {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	results := []interface{}{}
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

func (db *Imdb) Update(modelName, primaryKey string, keyValue, object interface{}) error {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	updatedCollection := []interface{}{}
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
		return fmt.Errorf("Unable to find a %v with %v = %v", modelName, primaryKey, keyValue)
	}
	db.Data[modelName] = updatedCollection
	return nil
}

func (db *Imdb) Delete(modelName, primaryKey string, keyValue interface{}) error {
	collection, ok := db.Data[modelName]
	if !ok {
		collection = []interface{}{}
	}
	updatedCollection := []interface{}{}
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
		return fmt.Errorf("Unable to find a %v with %v = %v", modelName, primaryKey, keyValue)
	}
	db.Data[modelName] = updatedCollection
	return nil
}
