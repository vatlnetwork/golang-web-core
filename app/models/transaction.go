package models

import (
	"fmt"
	"golang-web-core/app/domain"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (t TransactionModel) All() (any, error) {
	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
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

		transactions := []domain.Transaction{}
		err = cursor.All(ctx, &transactions)
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Create implements Model.
func (t TransactionModel) Create(object any) (any, error) {
	_, isTransaction := object.(domain.Transaction)
	if !isTransaction {
		return nil, fmt.Errorf("the given object is not a Transaction")
	}

	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		res, err := mongoAdapter.InsertOne(client, ctx, t.Name(), object)
		if err != nil {
			return nil, err
		}

		transaction := object.(domain.Transaction)
		transaction.Id = res.InsertedID.(bson.ObjectID)

		return transaction, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

// Delete implements Model.
func (t TransactionModel) Delete(key any) error {
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(key.(string))
		if err != nil {
			return err
		}

		filter := bson.M{
			t.PrimaryKey(): objectId,
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
func (t TransactionModel) DeleteWhere(query map[string]any) error {
	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter := bson.M{}
		for key := range query {
			filter[key] = query[key]
		}

		err = mongoAdapter.DeleteMany(client, ctx, t.Name(), filter)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrUnsupportedAdapter(t, t.adapter)
}

// Find implements Model.
func (t TransactionModel) Find(key any) (any, error) {
	_, isString := key.(string)
	if !isString {
		return nil, fmt.Errorf("key must be a string")
	}

	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(key.(string))
		if err != nil {
			return nil, err
		}

		filter := bson.M{
			t.PrimaryKey(): objectId,
		}

		cursor, err := mongoAdapter.Query(client, ctx, t.Name(), filter, nil)
		if err != nil {
			return nil, err
		}

		transactions := []domain.Transaction{}
		err = cursor.All(ctx, &transactions)
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
	return "transactions"
}

// PrimaryKey implements Model.
func (t TransactionModel) PrimaryKey() string {
	return "_id"
}

// Update implements Model.
func (t TransactionModel) Update(key any, object any) error {
	_, isString := key.(string)
	if !isString {
		return fmt.Errorf("key must be a string")
	}

	_, isTransaction := object.(domain.Transaction)
	if !isTransaction {
		return fmt.Errorf("the given object is not a Transaction")
	}

	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(key.(string))
		if err != nil {
			return err
		}

		filter := bson.M{
			t.PrimaryKey(): objectId,
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
func (t TransactionModel) UpdateWhere(query map[string]any, object any) error {
	_, isTransaction := object.(domain.Transaction)
	if !isTransaction {
		return fmt.Errorf("the given object is not a Transaction")
	}

	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
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
func (t TransactionModel) Where(query map[string]any) (any, error) {
	mongoAdapter, ok := (*t.adapter).(*mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, err
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		mongoQuery := bson.M{}
		for key := range query {
			mongoQuery[key] = query[key]
		}

		cursor, err := mongoAdapter.Query(client, ctx, t.Name(), mongoQuery, nil)
		if err != nil {
			return nil, err
		}

		transactions := []domain.Transaction{}
		err = cursor.All(ctx, &transactions)
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}

	return nil, ErrUnsupportedAdapter(t, t.adapter)
}

var _ Model = TransactionModel{}
