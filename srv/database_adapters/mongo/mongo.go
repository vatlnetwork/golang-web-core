package mongo

import (
	"context"
	"fmt"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/util"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.ConnectionString()))
	return client, ctx, cancel, err
}

func (m Mongo) Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}
