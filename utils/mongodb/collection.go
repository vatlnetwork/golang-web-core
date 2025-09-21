package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"golang-web-core/logging"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Collection struct {
	col    *mongo.Collection
	logger *logging.Logger
	ctx    context.Context
}

func NewCollection(col *mongo.Collection, logger *logging.Logger, ctx context.Context) (Collection, error) {
	if col == nil {
		return Collection{}, errors.New("collection is required")
	}

	if logger == nil {
		return Collection{}, errors.New("logger is required")
	}

	return Collection{
		col:    col,
		logger: logger,
		ctx:    ctx,
	}, nil
}

// DeleteMany
func (c Collection) DeleteMany(filter any) error {
	res, err := c.col.DeleteMany(c.ctx, filter)
	if err != nil {
		return err
	}

	c.logger.Debugf("Deleted %v documents from collection %v", res.DeletedCount, c.col.Name())

	return nil
}

// DeleteOne
func (c Collection) DeleteOne(filter any) error {
	res, err := c.col.DeleteOne(c.ctx, filter)
	if err != nil {
		return err
	}

	c.logger.Debugf("Deleted %v document from collection %v", res.DeletedCount, c.col.Name())

	return nil
}

// Find
func (c Collection) Find(filter any, result any) error {
	cursor, err := c.col.Find(c.ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(c.ctx)

	results := []map[string]any{}
	err = cursor.All(c.ctx, result)
	if err != nil {
		return err
	}

	c.logger.Debugf("Found %v documents in collection %v", len(results), c.col.Name())

	bytes, err := json.Marshal(results)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}

	return nil
}

// FindOne
func (c Collection) FindOne(filter any, result any) error {
	res := c.col.FindOne(c.ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}

	err := res.Decode(result)
	if err != nil {
		return err
	}

	c.logger.Debugf("Found 1 document in collection %v", c.col.Name())

	return nil
}

// InsertMany
func (c Collection) InsertMany(documents any) (*mongo.InsertManyResult, error) {
	res, err := c.col.InsertMany(c.ctx, documents)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf("Inserted %v documents into collection %v", len(res.InsertedIDs), c.col.Name())

	return res, nil
}

// InsertOne
func (c Collection) InsertOne(document any) (*mongo.InsertOneResult, error) {
	res, err := c.col.InsertOne(c.ctx, document)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf("Inserted 1 document into collection %v", c.col.Name())

	return res, nil
}

// UpdateMany
func (c Collection) UpdateMany(filter any, update any) error {
	res, err := c.col.UpdateMany(c.ctx, filter, update)
	if err != nil {
		return err
	}

	c.logger.Debugf("Updated %v documents in collection %v", res.ModifiedCount, c.col.Name())

	return nil
}

// UpdateOne
func (c Collection) UpdateOne(filter any, update any) error {
	res, err := c.col.UpdateOne(c.ctx, filter, update)
	if err != nil {
		return err
	}

	c.logger.Debugf("Updated %v document in collection %v", res.ModifiedCount, c.col.Name())

	return nil
}
