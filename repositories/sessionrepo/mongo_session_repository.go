package sessionrepo

import (
	"context"
	"errors"
	"fmt"
	"golang-web-core/domain"
	"golang-web-core/logging"
	"golang-web-core/utils/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoSessionRepoConfig struct {
	MongoConfig mongodb.Config `json:"mongoConfig"`
	Collection  string         `json:"collection"`
}

type MongoSessionRepository struct {
	config  MongoSessionRepoConfig
	logger  *logging.Logger
	adapter *mongodb.MongoAdapter
}

func NewMongoSessionRepository(config MongoSessionRepoConfig, logger *logging.Logger) (MongoSessionRepository, error) {
	if logger == nil {
		return MongoSessionRepository{}, errors.New("logger is required")
	}

	adapter, err := mongodb.NewMongoAdapter(config.MongoConfig, logger)
	if err != nil {
		return MongoSessionRepository{}, fmt.Errorf("failed to create mongodb adapter: %w", err)
	}

	connectionError := errors.New("connection not tested yet")
	for connectionError != nil {
		connectionError = adapter.TestConnection()
		if connectionError != nil {
			logger.Errorf("failed to connect to mongodb: %v", connectionError)
		}
	}

	return MongoSessionRepository{
		config:  config,
		logger:  logger,
		adapter: &adapter,
	}, nil
}

func (m *MongoSessionRepository) connect() (*mongo.Client, context.Context, context.CancelFunc, mongodb.Collection, error) {
	client, ctx, cancel, err := m.adapter.Connect()
	if err != nil {
		return nil, nil, nil, mongodb.Collection{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	collection, err := m.adapter.GetCollection(client, ctx, m.config.Collection)
	if err != nil {
		return nil, nil, nil, mongodb.Collection{}, fmt.Errorf("failed to get collection: %w", err)
	}

	return client, ctx, cancel, collection, nil
}

// CreateSession implements domain.SessionRepository.
func (m *MongoSessionRepository) CreateSession(session domain.Session) (domain.Session, error) {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoSession, err := MongoSessionFromDomain(session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to convert session to mongo session: %w", err)
	}

	res, err := collection.InsertOne(mongoSession)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to insert session: %w", err)
	}

	insertedId, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.Session{}, errors.New("inserted id is not a bson.ObjectID")
	}

	mongoSession.Id = insertedId

	return mongoSession.ToDomain(), nil
}

// DeleteAllForUser implements domain.SessionRepository.
func (m *MongoSessionRepository) DeleteAllForUser(userId string) error {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	return collection.DeleteMany(bson.M{"userId": userId})
}

// DeleteSession implements domain.SessionRepository.
func (m *MongoSessionRepository) DeleteSession(sessionId string) error {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(sessionId)
	if err != nil {
		return fmt.Errorf("failed to convert id to bson.ObjectID: %w", err)
	}

	return collection.DeleteOne(bson.M{"_id": mongoId})
}

// GetAllForUser implements domain.SessionRepository.
func (m *MongoSessionRepository) GetAllForUser(userId string) ([]domain.Session, error) {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoSessions := []MongoSession{}
	err = collection.Find(bson.M{"userId": userId}, &mongoSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions: %w", err)
	}

	sessions := []domain.Session{}
	for _, mongoSession := range mongoSessions {
		sessions = append(sessions, mongoSession.ToDomain())
	}

	return sessions, nil
}

// GetSession implements domain.SessionRepository.
func (m *MongoSessionRepository) GetSession(sessionId string) (domain.Session, error) {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(sessionId)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to convert id to bson.ObjectID: %w", err)
	}

	mongoSession := MongoSession{}
	err = collection.FindOne(bson.M{"_id": mongoId}, &mongoSession)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return domain.Session{}, errors.New(domain.ErrorSessionNotFound)
		}
		return domain.Session{}, fmt.Errorf("failed to find session: %w", err)
	}

	return mongoSession.ToDomain(), nil
}

var _ domain.SessionRepository = &MongoSessionRepository{}
