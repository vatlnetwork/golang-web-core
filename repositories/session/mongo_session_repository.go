package sessionrepo

import (
	"errors"
	"inventory-app/domain"
	"inventory-app/util/database_adapters/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const sessionCollection string = "sessions"

type MongoSessionRepository struct {
	mongoConfig          mongo.Config
	logMongoTransactions bool
}

func NewMongoSessionRepository(mongoConfig mongo.Config, logMongoTransactions bool) MongoSessionRepository {
	return MongoSessionRepository{
		mongoConfig:          mongoConfig,
		logMongoTransactions: logMongoTransactions,
	}
}

func (m MongoSessionRepository) adapter() *mongo.Mongo {
	return mongo.NewMongoAdapter(m.mongoConfig, m.logMongoTransactions)
}

// CreateSession implements domain.SessionRepository.
func (m MongoSessionRepository) CreateSession(session domain.Session) (domain.Session, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.Session{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoSession, err := MongoSessionFromDomain(session)
	if err != nil {
		return domain.Session{}, err
	}

	result, err := adapter.InsertOne(client, ctx, sessionCollection, mongoSession)
	if err != nil {
		return domain.Session{}, err
	}

	sessionId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.Session{}, errors.New("failed to get inserted session id")
	}

	mongoSession.Id = sessionId

	return mongoSession.ToDomain(), nil
}

// DeleteAllForUser implements domain.SessionRepository.
func (m MongoSessionRepository) DeleteAllForUser(userId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	err = adapter.DeleteMany(client, ctx, sessionCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSession implements domain.SessionRepository.
func (m MongoSessionRepository) DeleteSession(sessionId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoSessionId, err := bson.ObjectIDFromHex(sessionId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mongoSessionId}

	err = adapter.DeleteOne(client, ctx, sessionCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetAllForUser implements domain.SessionRepository.
func (m MongoSessionRepository) GetAllForUser(userId string) ([]domain.Session, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	cursor, err := adapter.Query(client, ctx, sessionCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoSessions := []MongoSession{}
	err = cursor.All(ctx, &mongoSessions)
	if err != nil {
		return nil, err
	}

	sessions := []domain.Session{}
	for _, mongoSession := range mongoSessions {
		sessions = append(sessions, mongoSession.ToDomain())
	}

	return sessions, nil
}

// GetSession implements domain.SessionRepository.
func (m MongoSessionRepository) GetSession(sessionId string) (domain.Session, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.Session{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoSessionId, err := bson.ObjectIDFromHex(sessionId)
	if err != nil {
		return domain.Session{}, err
	}

	filter := bson.M{"_id": mongoSessionId}

	cursor, err := adapter.Query(client, ctx, sessionCollection, filter, nil)
	if err != nil {
		return domain.Session{}, err
	}

	mongoSessions := []MongoSession{}
	err = cursor.All(ctx, &mongoSessions)
	if err != nil {
		return domain.Session{}, err
	}

	if len(mongoSessions) == 0 {
		return domain.Session{}, errors.New(domain.ErrorSessionNotFound)
	}

	return mongoSessions[0].ToDomain(), nil
}

var _ domain.SessionRepository = MongoSessionRepository{}
