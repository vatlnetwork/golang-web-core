package moneylocationrepo

import (
	"errors"
	"inventory-app/domain"
	"inventory-app/util/database_adapters/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const moneyLocationCollection string = "moneyLocations"

type MongoMoneyLocationRepository struct {
	connectionConfig mongo.Config
	logTransactions  bool
}

func NewMongoMoneyLocationRepository(
	connectionConfig mongo.Config,
	logTransactions bool,
) MongoMoneyLocationRepository {
	return MongoMoneyLocationRepository{
		connectionConfig: connectionConfig,
		logTransactions:  logTransactions,
	}
}

func (m MongoMoneyLocationRepository) adapter() *mongo.Mongo {
	return mongo.NewMongoAdapter(m.connectionConfig, m.logTransactions)
}

// DeleteAllMoneyLocationsForUser implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) DeleteAllMoneyLocationsForUser(userId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	err = adapter.DeleteMany(client, ctx, moneyLocationCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// CreateMoneyLocation implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) CreateMoneyLocation(moneyLocation domain.MoneyLocation) (domain.MoneyLocation, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.MoneyLocation{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoMoneyLocation, err := MongoMoneyLocationFromDomain(moneyLocation)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	result, err := adapter.InsertOne(client, ctx, moneyLocationCollection, mongoMoneyLocation)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	insertedId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.MoneyLocation{}, errors.New("failed to get inserted ID")
	}

	mongoMoneyLocation.Id = insertedId

	return mongoMoneyLocation.ToDomain(), nil
}

// DeleteMoneyLocation implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) DeleteMoneyLocation(id string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId}

	err = adapter.DeleteOne(client, ctx, moneyLocationCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetMoneyLocation implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) GetMoneyLocation(id string) (domain.MoneyLocation, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.MoneyLocation{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	filter := bson.M{"_id": objectId}

	cursor, err := adapter.Query(client, ctx, moneyLocationCollection, filter, nil)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	mongoMoneyLocations := []MongoMoneyLocation{}
	err = cursor.All(ctx, &mongoMoneyLocations)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	if len(mongoMoneyLocations) == 0 {
		return domain.MoneyLocation{}, errors.New(domain.ErrorMoneyLocationNotFound)
	}

	return mongoMoneyLocations[0].ToDomain(), nil
}

// GetMoneyLocationsForUser implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) GetMoneyLocationsForUser(userId string) ([]domain.MoneyLocation, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	cursor, err := adapter.Query(client, ctx, moneyLocationCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoMoneyLocations := []MongoMoneyLocation{}
	err = cursor.All(ctx, &mongoMoneyLocations)
	if err != nil {
		return nil, err
	}

	moneyLocations := make([]domain.MoneyLocation, len(mongoMoneyLocations))
	for i, mongoMoneyLocation := range mongoMoneyLocations {
		moneyLocations[i] = mongoMoneyLocation.ToDomain()
	}

	return moneyLocations, nil
}

// UpdateMoneyLocation implements domain.MoneyLocationRepository.
func (m MongoMoneyLocationRepository) UpdateMoneyLocation(moneyLocation domain.MoneyLocation) (domain.MoneyLocation, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.MoneyLocation{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoMoneyLocation, err := MongoMoneyLocationFromDomain(moneyLocation)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	filter := bson.M{"_id": mongoMoneyLocation.Id}

	update := bson.M{"$set": mongoMoneyLocation}

	err = adapter.UpdateOne(client, ctx, moneyLocationCollection, filter, update)
	if err != nil {
		return domain.MoneyLocation{}, err
	}

	return mongoMoneyLocation.ToDomain(), nil
}

var _ domain.MoneyLocationRepository = MongoMoneyLocationRepository{}
