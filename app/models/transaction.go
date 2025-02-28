package models

import (
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/mongo"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	Id          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Description string  `json:"description"`
	GroupId     string  `json:"groupId"`
}

func NewTransaction(id, description, groupId string, amount float64, timestamp int64) Transaction {
	Id := id
	if Id == "" {
		Id = uuid.NewString()
	}
	Timestamp := timestamp
	if Timestamp == 0 {
		Timestamp = time.Now().UnixMilli()
	}

	return Transaction{
		Id:          Id,
		Amount:      amount,
		Timestamp:   Timestamp,
		Description: description,
		GroupId:     groupId,
	}
}

type TransactionModel struct {
	adapter *databaseadapters.DatabaseAdapter
}

func NewTransactionModel(adapter *databaseadapters.DatabaseAdapter) TransactionModel {
	return TransactionModel{
		adapter: adapter,
	}
}

// Adapter implements Model.
func (t TransactionModel) Adapter() *databaseadapters.DatabaseAdapter {
	return t.adapter
}

// All implements Model.
func (t TransactionModel) All() (interface{}, error) {
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

		transactions := []Transaction{}
		err = cursor.All(ctx, &transactions)
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Create implements Model.
func (t TransactionModel) Create(object interface{}) (interface{}, error) {
	_, isObject := object.(Transaction)
	if !isObject {
		return nil, fmt.Errorf("the given object is not a Transaction")
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
func (t TransactionModel) Delete(key interface{}) error {
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
func (t TransactionModel) DeleteWhere(query map[string]interface{}) error {
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
func (t TransactionModel) Find(key interface{}) (interface{}, error) {
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

		transactions := []Transaction{}
		err = cursor.All(context, &transactions)
		if err != nil {
			return nil, err
		}

		if len(transactions) == 0 {
			return nil, fmt.Errorf("unable to find a Transaction with %v = %v", t.PrimaryKey(), key)
		}

		return transactions[0], nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Name implements Model.
func (t TransactionModel) Name() string {
	// this is the name of the collection or table the data for this model is stored in
	return "transactions"
}

// PrimaryKey implements Model.
func (t TransactionModel) PrimaryKey() string {
	return "Id"
}

// Update implements Model.
func (t TransactionModel) Update(key interface{}, object interface{}) error {
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	_, isTransaction := object.(Transaction)
	if !isTransaction {
		return fmt.Errorf("object must be a Transaction")
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
func (t TransactionModel) UpdateWhere(query map[string]interface{}, object interface{}) error {
	_, isTransaction := object.(Transaction)
	if !isTransaction {
		return fmt.Errorf("object must be a Transaction")
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
func (t TransactionModel) Where(query map[string]interface{}) (interface{}, error) {
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

		transactions := []Transaction{}
		err = cursor.All(context, &transactions)
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

var TransactionModelVerifier Model = TransactionModel{}
