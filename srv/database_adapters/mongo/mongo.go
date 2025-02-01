package mongo

import (
	"context"
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/util"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	databaseadapters.ConnectionConfig
}

func NewMongoAdapter() Mongo {
	return Mongo{
		ConnectionConfig: DefaultConfig(),
	}
}

func (m Mongo) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m Mongo) Connection() databaseadapters.ConnectionConfig {
	return m.ConnectionConfig
}

func (m Mongo) ConnectionString() string {
	authAndHost := m.Hostname
	if m.UsingAuth() {
		authAndHost = fmt.Sprintf("%v:%v@%v", m.Username, m.Password, m.Hostname)
	}
	return fmt.Sprintf("mongodb://%v/%v", authAndHost, m.Database)
}

func (m Mongo) TestConnection() error {
	client, context, cancel, err := m.Connect()
	if err != nil {
		return err
	}
	defer m.Close(client, context, cancel)
	return m.Ping(client, context)
}

func (m Mongo) Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			util.LogFatal(err)
		}
	}()
}

func (m Mongo) Connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(options.Client().ApplyURI(m.ConnectionString()))
	return client, ctx, cancel, err
}

func (m Mongo) Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m Mongo) InsertOne(client *mongo.Client, ctx context.Context, col string, doc interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.InsertOne(ctx, doc)
	return err
}

func (m Mongo) InsertMany(client *mongo.Client, ctx context.Context, col string, docs []interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.InsertMany(ctx, docs)
	return err
}

func (m Mongo) Query(client *mongo.Client, ctx context.Context, col string, query, field interface{}) (*mongo.Cursor, error) {
	collection := client.Database(m.Database).Collection(col)
	result, err := collection.Find(ctx, query, options.Find().SetProjection(field))
	return result, err
}

func (m Mongo) UpdateOne(client *mongo.Client, ctx context.Context, col string, filter, update interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (m Mongo) UpdateMany(client *mongo.Client, ctx context.Context, col string, filter, update interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.UpdateMany(ctx, filter, update)
	return err
}

func (m Mongo) DeleteOne(client *mongo.Client, ctx context.Context, col string, query interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.DeleteOne(ctx, query)
	return err
}

func (m Mongo) DeleteMany(client *mongo.Client, ctx context.Context, col string, query interface{}) error {
	collection := client.Database(m.Database).Collection(col)
	_, err := collection.DeleteMany(ctx, query)
	return err
}
