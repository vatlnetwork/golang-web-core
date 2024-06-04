package mongo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri+"?authSource=admin"))
	return client, ctx, cancel, err
}

func Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase string, col string, doc interface{}) error {
	collection := client.Database(dataBase).Collection(col)
	_, err := collection.InsertOne(ctx, doc)
	c := color.New(color.FgGreen)
	c.Printf("Insert into %v/%v: %v\n", dataBase, col, doc)
	return err
}

func InsertMany(client *mongo.Client, ctx context.Context, dataBase string, col string, docs []interface{}) error {
	collection := client.Database(dataBase).Collection(col)
	_, err := collection.InsertMany(ctx, docs)
	c := color.New(color.FgGreen)
	c.Printf("Insert %v records into %v/%v: %v\n", len(docs), dataBase, col, docs)
	return err
}

// pass 0 into filter to get all of the matching results \n
func Query(client *mongo.Client, ctx context.Context, dataBase string, col string, query interface{}, pointer interface{}, limit int64) error {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.Find(ctx, query, options.Find().SetLimit(limit))
	if err != nil {
		return err
	}
	err = result.All(ctx, pointer)
	if err != nil {
		return err
	}
	pointerVal := reflect.ValueOf(pointer).Elem()
	c := color.New(color.FgBlue)
	c.Printf("Returned %v results from %v/%v: %v\n", pointerVal.Len(), dataBase, col, query)
	return nil
}

func AggregatedQuery(client *mongo.Client, ctx context.Context, dataBase string, col string, agg bson.A, pointer interface{}) error {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.Aggregate(ctx, agg)
	if err != nil {
		return err
	}
	err = result.All(ctx, pointer)
	if err != nil {
		return err
	}
	pointerVal := reflect.ValueOf(pointer).Elem()
	c := color.New(color.FgBlue)
	c.Printf("Returned %v results from %v/%v: %v\n", pointerVal.Len(), dataBase, col, agg)
	return nil
}

func UpdateOne(
	client *mongo.Client,
	ctx context.Context,
	dataBase string,
	col string,
	filter interface{},
	update interface{},
) error {
	collection := client.Database(dataBase).Collection(col)
	_, err := collection.UpdateOne(ctx, filter, update)
	c := color.New(color.FgYellow)
	c.Printf("Updated 1 record in %v/%v: %v\n", dataBase, col, update)
	return err
}

func UpdateMany(
	client *mongo.Client,
	ctx context.Context,
	dataBase string,
	col string,
	filter interface{},
	update interface{},
) error {
	collection := client.Database(dataBase).Collection(col)
	res, err := collection.UpdateMany(ctx, filter, update)
	c := color.New(color.FgYellow)
	c.Printf("Updated %v records in %v/%v: %v\n", res.MatchedCount, dataBase, col, update)
	return err
}

func DeleteOne(client *mongo.Client, ctx context.Context, dataBase string, col string, query interface{}) error {
	collection := client.Database(dataBase).Collection(col)
	_, err := collection.DeleteOne(ctx, query)
	c := color.New(color.FgRed)
	c.Printf("Deleted 1 record in %v/%v: %v\n", dataBase, col, query)
	return err
}

func DeleteMany(client *mongo.Client, ctx context.Context, dataBase string, col string, query interface{}) error {
	collection := client.Database(dataBase).Collection(col)
	res, err := collection.DeleteMany(ctx, query)
	c := color.New(color.FgRed)
	c.Printf("Deleted %v records in %v/%v: %v\n", res.DeletedCount, dataBase, col, query)
	return err
}
