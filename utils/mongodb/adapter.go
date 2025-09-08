package mongodb

import (
	"context"
	"errors"
	"golang-web-core/logging"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoAdapter struct {
	config Config
	logger *logging.Logger
}

func NewMongoAdapter(config Config, logger *logging.Logger) (MongoAdapter, error) {
	if !config.IsEnabled() {
		return MongoAdapter{}, errors.New("mongodb is missing hostname or database")
	}

	if config.TimeoutSeconds <= 0 {
		return MongoAdapter{}, errors.New("timeout seconds must be greater than 0")
	}

	if logger == nil {
		return MongoAdapter{}, errors.New("logger is required")
	}

	return MongoAdapter{
		config: config,
		logger: logger,
	}, nil
}

func (m MongoAdapter) TestConnection() error {
	client, ctx, cancel, err := m.Connect()
	if err != nil {
		return err
	}
	err = m.Ping(client, ctx)
	if err != nil {
		return err
	}
	m.Close(client, ctx, cancel)
	return nil
}

func (m MongoAdapter) Connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.TimeoutSeconds)*time.Second)
	client, err := mongo.Connect(options.Client().ApplyURI(m.config.ConnectionString()))
	return client, ctx, cancel, err
}

func (m MongoAdapter) Ping(client *mongo.Client, ctx context.Context) error {
	return client.Ping(ctx, readpref.Primary())
}

func (m MongoAdapter) Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	err := client.Disconnect(ctx)
	if err != nil {
		m.logger.Errorf("Failed to disconnect from mongodb: %v", err)
	}
	cancel()
}

func (m MongoAdapter) GetCollection(client *mongo.Client, context context.Context, collectionName string) (Collection, error) {
	return NewCollection(client.Database(m.config.Database).Collection(collectionName), m.logger, context)
}
