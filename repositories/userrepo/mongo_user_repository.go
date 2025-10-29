package userrepo

import (
	"context"
	"errors"
	"fmt"
	"golang-web-core/domain"
	"golang-web-core/logging"
	"golang-web-core/utils/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepoConfig struct {
	MongoConfig mongodb.Config `json:"mongoConfig"`
	Collection  string         `json:"collection"`
}

type MongoUserRepository struct {
	config  MongoUserRepoConfig
	logger  *logging.Logger
	adapter *mongodb.MongoAdapter
}

func NewMongoUserRepository(config MongoUserRepoConfig, logger *logging.Logger) (MongoUserRepository, error) {
	if logger == nil {
		return MongoUserRepository{}, errors.New("logger is required")
	}

	adapter, err := mongodb.NewMongoAdapter(config.MongoConfig, logger)
	if err != nil {
		return MongoUserRepository{}, fmt.Errorf("failed to create mongodb adapter: %w", err)
	}

	return MongoUserRepository{
		config:  config,
		logger:  logger,
		adapter: &adapter,
	}, nil
}

func (m *MongoUserRepository) connect() (*mongo.Client, context.Context, context.CancelFunc, mongodb.Collection, error) {
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

// CreateUser implements domain.UserRepository.
func (m *MongoUserRepository) CreateUser(user domain.User) (domain.User, error) {
	userWithEmail, err := m.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() != domain.ErrorUserNotFound {
			return domain.User{}, err
		}
	}
	if userWithEmail.Id != "" {
		return domain.User{}, errors.New(domain.ErrorUserAlreadyExists)
	}

	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoUser, err := MongoUserFromDomain(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to convert user to mongo user: %w", err)
	}

	res, err := collection.InsertOne(mongoUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	insertedId, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.User{}, errors.New("inserted id is not a bson.ObjectID")
	}

	mongoUser.Id = insertedId

	return mongoUser.ToDomain(), nil
}

// DeleteUser implements domain.UserRepository.
func (m *MongoUserRepository) DeleteUser(userId string) error {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("failed to convert id to bson.ObjectID: %w", err)
	}

	return collection.DeleteOne(bson.M{"_id": mongoId})
}

// GetUser implements domain.UserRepository.
func (m *MongoUserRepository) GetUser(userId string) (domain.User, error) {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to convert id to bson.ObjectID: %w", err)
	}

	mongoUser := MongoUser{}
	err = collection.FindOne(bson.M{"_id": mongoId}, &mongoUser)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return domain.User{}, errors.New(domain.ErrorUserNotFound)
		}
		return domain.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return mongoUser.ToDomain(), nil
}

// GetUserByEmail implements domain.UserRepository.
func (m *MongoUserRepository) GetUserByEmail(email string) (domain.User, error) {
	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	mongoUser := MongoUser{}
	err = collection.FindOne(bson.M{"email": email}, &mongoUser)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return domain.User{}, errors.New(domain.ErrorUserNotFound)
		}
		return domain.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return mongoUser.ToDomain(), nil
}

// UpdateUser implements domain.UserRepository.
func (m *MongoUserRepository) UpdateUser(user domain.User) (domain.User, error) {
	originalUser, err := m.GetUser(user.Id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get original user: %w", err)
	}

	if originalUser.Email != user.Email {
		userWithEmail, err := m.GetUserByEmail(user.Email)
		if err != nil {
			if err.Error() != domain.ErrorUserNotFound {
				return domain.User{}, fmt.Errorf("failed to get user with email: %w", err)
			}
		}
		if userWithEmail.Id != "" {
			return domain.User{}, errors.New(domain.ErrorUserAlreadyExists)
		}
	}

	client, ctx, cancel, collection, err := m.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer m.adapter.Close(client, ctx, cancel)

	user.UpdatedAt = time.Now()

	mongoUser, err := MongoUserFromDomain(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to convert user to mongo user: %w", err)
	}

	filter := bson.M{"_id": mongoUser.Id}
	update := bson.M{"$set": mongoUser}

	err = collection.UpdateOne(filter, update)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return mongoUser.ToDomain(), nil
}

var _ domain.UserRepository = &MongoUserRepository{}
