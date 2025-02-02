package models

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/imdb"
	"golang-web-core/srv/database_adapters/mongo"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var UnknownAdapterError error = fmt.Errorf("The database adapter specified for TestModel is unrecognized.")

// this line is here to verify that TestModel implements the Model interface
var TestModelVerifier Model = TestModel{}

type TestObject struct {
	Id      string `json:"id" bson:"id"`
	Number  int    `json:"number" bson:"number"`
	Boolean bool   `json:"boolean" bson:"boolean"`
}

func NewTestObject(number int, boolean bool) TestObject {
	return TestObject{
		Id:      uuid.NewString(),
		Number:  number,
		Boolean: boolean,
	}
}

type TestModel struct {
	adapter *databaseadapters.DatabaseAdapter
}

func NewTestModel(adapter *databaseadapters.DatabaseAdapter) TestModel {
	return TestModel{
		adapter: adapter,
	}
}

func (m TestModel) Adapter() *databaseadapters.DatabaseAdapter {
	return m.adapter
}

func (m TestModel) Name() string {
	return "testObjects"
}

func (m TestModel) PrimaryKey() string {
	return "id"
}

func (m TestModel) Create(object interface{}) (interface{}, error) {
	// verify that object is a TestObject
	_, isObject := object.(TestObject)
	if !isObject {
		return nil, fmt.Errorf("the given object is not a TestObject")
	}

	// if the adapter is a mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to the mongo database
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		// insert the record
		err = mongoAdapter.InsertOne(client, context, m.Name(), object)
		if err != nil {
			return nil, err
		}

		// return the object
		return object, nil
	}

	// if the adapter is an imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// insert the object
		imdbAdapter.Insert(m.Name(), object)

		// return the object
		return object, nil
	}

	// this line is only reached if the function fails to return before this,
	// in which case it is assumed that the adapter that was passed in was not
	// recognized as a known adapter
	return nil, UnknownAdapterError
}

func (m TestModel) Find(key interface{}) (interface{}, error) {
	// make sure key is a string
	_, isString := key.(string)
	if !isString {
		return nil, fmt.Errorf("key must be a string")
	}

	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to mongo database
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		// build query
		query := bson.M{
			m.PrimaryKey(): key,
		}

		// get results from mongo database
		cursor, err := mongoAdapter.Query(client, context, m.Name(), query, nil)
		if err != nil {
			return nil, err
		}

		// decode the results from the mongo database
		objects := []TestObject{}
		err = cursor.All(context, &objects)
		if err != nil {
			return nil, err
		}

		// make sure there is a first object
		if len(objects) == 0 {
			return nil, fmt.Errorf("unable to find a TestObject with %v = %v", m.PrimaryKey(), key)
		}

		// return the first object
		return objects[0], nil
	}

	// case for imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// get the object from the memory collection
		iface, err := imdbAdapter.Find(m.Name(), m.PrimaryKey(), key)
		if err != nil {
			return nil, err
		}

		// verify the object is a TestObject
		object, ok := iface.(TestObject)
		if !ok {
			return nil, fmt.Errorf("the returned object was not a TestObject")
		}

		// return the object
		return object, nil
	}

	// return unknown adapter error if the adapter passed in can not be matched to any of the adapter cases
	return nil, UnknownAdapterError
}

func (m TestModel) Where(query map[string]interface{}) (interface{}, error) {
	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to mongo db
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		// build query
		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		// get data from the db
		cursor, err := mongoAdapter.Query(client, context, m.Name(), mongoQuery, nil)
		if err != nil {
			return nil, err
		}

		// decode the data from the db
		objects := []TestObject{}
		err = cursor.All(context, &objects)
		if err != nil {
			return nil, err
		}

		// return the decoded objects
		return objects, nil
	}

	// case for imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// get the records from the memory collection
		records := imdbAdapter.Query(m.Name(), query)

		// convert the records into TestObjects
		results := []TestObject{}
		for _, record := range records {
			// verify the record is a TestObject
			obj, ok := record.(TestObject)
			if !ok {
				return nil, fmt.Errorf("found a record in the results that isn't a TestObject")
			}

			// append the TestObject to the results array
			results = append(results, obj)
		}

		// return the results
		return results, nil
	}

	// return unknown adapter error if the passed in adapter is not matched to a case
	return nil, UnknownAdapterError
}

func (m TestModel) All() (interface{}, error) {
	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to the mongo db
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// get all records from the database
		cursor, err := mongoAdapter.Query(client, ctx, m.Name(), bson.M{}, nil)
		if err != nil {
			return nil, err
		}

		// decode all of the records
		var objects []TestObject
		err = cursor.All(ctx, &objects)
		if err != nil {
			return nil, err
		}

		// return the decoded records
		return objects, nil
	}

	// case for imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// get all of the objects from the memory collection
		objects := imdbAdapter.GetAll(m.Name())

		// convert the results into TestObjects
		results := []TestObject{}
		for _, object := range objects {
			_, isObject := object.(TestObject)
			if !isObject {
				return nil, fmt.Errorf("found an object that is not a TestObject")
			}
			results = append(results, object.(TestObject))
		}

		// return the results
		return results, nil
	}

	// return unknown adapter error if the passed in adapter is not matched to a case
	return nil, UnknownAdapterError
}

func (m TestModel) Update(key, object interface{}) error {
	// verify key is a string
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	// verify object is a TestObject
	_, isObject := object.(TestObject)
	if !isObject {
		return fmt.Errorf("object must be a TestObject")
	}

	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to mongo db
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// build query
		filter := bson.M{
			m.PrimaryKey(): key,
		}

		// build update object
		update := bson.M{
			"$set": object,
		}

		// update the db
		err = mongoAdapter.UpdateOne(client, ctx, m.Name(), filter, update)
		if err != nil {
			return err
		}

		// return nil to indicate the operation succeeded without error
		return nil
	}

	// case for imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// update the memory collection
		err := imdbAdapter.Update(m.Name(), m.PrimaryKey(), key, object)
		if err != nil {
			return err
		}

		// return nil to indicate absence of errors
		return nil
	}

	// return unknown adapter error unless adapter is matched to a case
	return UnknownAdapterError
}

func (m TestModel) UpdateWhere(query map[string]interface{}, object interface{}) error {
	// verify object is a TestObject
	_, isObject := object.(TestObject)
	if !isObject {
		return fmt.Errorf("object must be a TestObject")
	}

	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to the mongo db
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// build the update query that will be used to identify objects that will be replaced
		// by the update object
		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		// build the update object
		update := bson.M{
			"$set": object,
		}

		// update the db
		err = mongoAdapter.UpdateMany(client, ctx, m.Name(), mongoQuery, update)
		if err != nil {
			return err
		}

		// return nil to indicate absence of errors
		return nil
	}

	// case for imdb adapter
	_, ok = (*m.adapter).(imdb.Imdb)
	if ok {
		// get records matching query
		records, err := m.Where(query)
		if err != nil {
			return err
		}
		objects := records.([]TestObject)

		// update each object matching the query with the new object
		for _, obj := range objects {
			err := m.Update(obj.Id, object)
			if err != nil {
				return err
			}
		}

		// return nil to indicate absence of errors
		return nil
	}

	// return unknown adapter error for adapters without a case
	return UnknownAdapterError
}

func (m TestModel) Delete(key interface{}) error {
	// verify key is a string
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect the mongo db
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// build query
		filter := bson.M{
			m.PrimaryKey(): key,
		}

		// delete the first record matching the query
		err = mongoAdapter.DeleteOne(client, ctx, m.Name(), filter)
		if err != nil {
			return err
		}

		// return nil to indicate the absence of errors
		return nil
	}

	// case for imdb adapter
	imdbAdapter, ok := (*m.adapter).(imdb.Imdb)
	if ok {
		// delete the record matching the key
		err := imdbAdapter.Delete(m.Name(), m.PrimaryKey(), key)
		if err != nil {
			return err
		}

		// return nil to indicate absence of errors
		return nil
	}

	// return unknown adapter error for adapters without a case
	return UnknownAdapterError
}

func (m TestModel) DeleteWhere(query map[string]interface{}) error {
	// case for mongo adapter
	mongoAdapter, ok := (*m.adapter).(mongo.Mongo)
	if ok {
		// connect to the mongo db
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// build the query
		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		// delete all records matched by the query
		err = mongoAdapter.DeleteMany(client, ctx, m.Name(), mongoQuery)
		if err != nil {
			return err
		}

		// return nil to indicate absence of errors
		return nil
	}

	// case for imdb adapter
	_, ok = (*m.adapter).(imdb.Imdb)
	if ok {
		// find all records matching the query
		records, err := m.Where(query)
		if err != nil {
			return err
		}
		objects := records.([]TestObject)

		// delete each object from the memory collection
		for _, object := range objects {
			err := m.Delete(object.Id)
			if err != nil {
				return err
			}
		}

		// return nil to indicate the absence of errors
		return nil
	}

	// return unknown adapter error for adapters that don't have a case written here
	return UnknownAdapterError
}
