package models

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/mongo"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionGroup struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

func NewTransactionGroup(id, description string) TransactionGroup {
	Id := id
	if Id == "" {
		Id = uuid.NewString()
	}
	Description := description
	if Description == "" {
		Description = "Unnamed Group"
	}

	return TransactionGroup{
		Id:          Id,
		Description: Description,
	}
}

type TransactionGroupModel struct {
	adapter *databaseadapters.DatabaseAdapter
}

func NewTransactionGroupModel(adapter *databaseadapters.DatabaseAdapter) TransactionGroupModel {
	return TransactionGroupModel{
		adapter: adapter,
	}
}

// Adapter implements Model.
func (t TransactionGroupModel) Adapter() *databaseadapters.DatabaseAdapter {
	return t.adapter
}

// All implements Model.
func (t TransactionGroupModel) All() (interface{}, error) {
	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		cursor, err := mongoAdapter.Query(client, ctx, t.Name(), bson.M{}, nil)
		if err != nil {
			return nil, err
		}

		groups := []TransactionGroup{}
		err = cursor.All(ctx, &groups)
		if err != nil {
			return nil, err
		}

		return groups, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Create implements Model.
func (t TransactionGroupModel) Create(object interface{}) (interface{}, error) {
	_, isObject := object.(TransactionGroup)
	if !isObject {
		return nil, fmt.Errorf("the given object is not a TransactionGroup")
	}

	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		err = mongoAdapter.InsertOne(client, context, t.Name(), object)
		if err != nil {
			return nil, err
		}

		return object, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Delete implements Model.
func (t TransactionGroupModel) Delete(key interface{}) error {
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		lowerCaseKey := strings.ToLower(t.PrimaryKey())
		filter := bson.M{
			lowerCaseKey: key,
		}

		err = mongoAdapter.DeleteOne(client, ctx, t.Name(), filter)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrUnsupportedAdapter(t, t.adapter)
}

// DeleteWhere implements Model.
func (t TransactionGroupModel) DeleteWhere(query map[string]interface{}) error {
	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		err = mongoAdapter.DeleteMany(client, ctx, t.Name(), mongoQuery)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrUnsupportedAdapter(t, t.adapter)
}

// Find implements Model.
func (t TransactionGroupModel) Find(key interface{}) (interface{}, error) {
	_, isString := key.(string)
	if !isString {
		return nil, fmt.Errorf("key must be a string")
	}

	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		lowerCaseKey := strings.ToLower(t.PrimaryKey())
		query := bson.M{
			lowerCaseKey: key,
		}

		cursor, err := mongoAdapter.Query(client, context, t.Name(), query, nil)
		if err != nil {
			return nil, err
		}

		groups := []TransactionGroup{}
		err = cursor.All(context, &groups)
		if err != nil {
			return nil, err
		}

		if len(groups) == 0 {
			return nil, fmt.Errorf("unable to find a TransactionGroup with %v = %v", t.PrimaryKey(), key)
		}

		return groups[0], nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Name implements Model.
func (t TransactionGroupModel) Name() string {
	// this is the name of the collection or table the data for this model is stored in
	return "transactionGroups"
}

// PrimaryKey implements Model.
func (t TransactionGroupModel) PrimaryKey() string {
	return "Id"
}

// Update implements Model.
func (t TransactionGroupModel) Update(key interface{}, object interface{}) error {
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	_, isTransactionGroup := object.(TransactionGroup)
	if !isTransactionGroup {
		return fmt.Errorf("object must be a TransactionGroup")
	}

	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		lowerCaseKey := strings.ToLower(t.PrimaryKey())
		filter := bson.M{
			lowerCaseKey: key,
		}

		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateOne(client, ctx, t.Name(), filter, update)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrUnsupportedAdapter(t, t.adapter)
}

// UpdateWhere implements Model.
func (t TransactionGroupModel) UpdateWhere(query map[string]interface{}, object interface{}) error {
	_, isTransactionGroup := object.(TransactionGroup)
	if !isTransactionGroup {
		return fmt.Errorf("object must be a TransactionGroup")
	}

	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateMany(client, ctx, t.Name(), mongoQuery, update)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrUnsupportedAdapter(t, t.adapter)
}

// Where implements Model.
func (t TransactionGroupModel) Where(query map[string]interface{}) (interface{}, error) {
	mongoAdapter, ok := (*t.adapter).(mongo.Mongo)
	if ok {
		client, context, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, context, cancel)

		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		cursor, err := mongoAdapter.Query(client, context, t.Name(), mongoQuery, nil)
		if err != nil {
			return nil, err
		}

		groups := []TransactionGroup{}
		err = cursor.All(context, &groups)
		if err != nil {
			return nil, err
		}

		return groups, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

var TransactionGroupModelVerifier Model = TransactionGroupModel{}
