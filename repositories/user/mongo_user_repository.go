package userrepo

import (
	"errors"
	"inventory-app/domain"
	"inventory-app/util/database_adapters/mongo"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const userCollection string = "users"

type MongoUserRepository struct {
	mongoConfig          mongo.Config
	logMongoTransactions bool
}

func NewMongoUserRepository(mongoConfig mongo.Config, logMongoTransactions bool) MongoUserRepository {
	return MongoUserRepository{
		mongoConfig:          mongoConfig,
		logMongoTransactions: logMongoTransactions,
	}
}

func (m MongoUserRepository) adapter() *mongo.Mongo {
	return mongo.NewMongoAdapter(m.mongoConfig, m.logMongoTransactions)
}

// CreateUser implements domain.UserRepository.
func (m MongoUserRepository) CreateUser(user domain.User) (domain.User, error) {
	matchingUser, err := m.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() != domain.ErrorUserNotFound {
			return domain.User{}, err
		}
	}
	if matchingUser.Id != "" {
		return domain.User{}, errors.New(domain.ErrorUserAlreadyExists)
	}

	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.User{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoUser, err := MongoUserFromDomain(user)
	if err != nil {
		return domain.User{}, err
	}

	result, err := adapter.InsertOne(client, ctx, userCollection, mongoUser)
	if err != nil {
		return domain.User{}, err
	}

	insertedId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.User{}, errors.New("failed to get inserted user id")
	}

	mongoUser.Id = insertedId

	return mongoUser.ToDomain(), nil
}

// DeleteUser implements domain.UserRepository.
func (m MongoUserRepository) DeleteUser(userId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoUserId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mongoUserId}

	err = adapter.DeleteOne(client, ctx, userCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetUser implements domain.UserRepository.
func (m MongoUserRepository) GetUser(userId string) (domain.User, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.User{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoUserId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"_id": mongoUserId}

	cursor, err := adapter.Query(client, ctx, userCollection, filter, nil)
	if err != nil {
		return domain.User{}, err
	}

	mongoUsers := []MongoUser{}
	err = cursor.All(ctx, &mongoUsers)
	if err != nil {
		return domain.User{}, err
	}

	if len(mongoUsers) == 0 {
		return domain.User{}, errors.New(domain.ErrorUserNotFound)
	}

	return mongoUsers[0].ToDomain(), nil
}

// GetUserByEmail implements domain.UserRepository.
func (m MongoUserRepository) GetUserByEmail(email string) (domain.User, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.User{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"email": email}

	cursor, err := adapter.Query(client, ctx, userCollection, filter, nil)
	if err != nil {
		return domain.User{}, err
	}

	mongoUsers := []MongoUser{}
	err = cursor.All(ctx, &mongoUsers)
	if err != nil {
		return domain.User{}, err
	}

	if len(mongoUsers) == 0 {
		return domain.User{}, errors.New(domain.ErrorUserNotFound)
	}

	return mongoUsers[0].ToDomain(), nil
}

// UpdateUser implements domain.UserRepository.
func (m MongoUserRepository) UpdateUser(user domain.User) (domain.User, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.User{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	user.UpdatedAt = time.Now()

	mongoUser, err := MongoUserFromDomain(user)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"_id": mongoUser.Id}

	update := bson.M{"$set": mongoUser}

	err = adapter.UpdateOne(client, ctx, userCollection, filter, update)
	if err != nil {
		return domain.User{}, err
	}

	return mongoUser.ToDomain(), nil
}

var _ domain.UserRepository = MongoUserRepository{}
