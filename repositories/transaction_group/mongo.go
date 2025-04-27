package transactiongrouprepo

import (
	"errors"
	"inventory-app/domain"
	"inventory-app/util/database_adapters/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const transactionGroupCollection string = "transactionGroups"

var ErrorTransactionGroupNotFound error = errors.New("transaction group not found")

type MongoTransactionGroupRepository struct {
	connectionConfig mongo.Config
	logTransactions  bool
}

func NewMongoTransactionGroupRepository(connectionConfig mongo.Config, logTransactions bool) MongoTransactionGroupRepository {
	return MongoTransactionGroupRepository{
		connectionConfig: connectionConfig,
		logTransactions:  logTransactions,
	}
}

func (m MongoTransactionGroupRepository) adapter() *mongo.Mongo {
	return mongo.NewMongoAdapter(m.connectionConfig, m.logTransactions)
}

// CreateTransactionGroup implements domain.TransactionGroupRepository.
func (m MongoTransactionGroupRepository) CreateTransactionGroup(transactionGroup domain.TransactionGroup) (domain.TransactionGroup, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.TransactionGroup{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoTransactionGroup, err := MongoTransactionGroupFromDomain(transactionGroup)
	if err != nil {
		return domain.TransactionGroup{}, err
	}

	result, err := adapter.InsertOne(client, ctx, transactionGroupCollection, mongoTransactionGroup)
	if err != nil {
		return domain.TransactionGroup{}, err
	}

	insertedId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.TransactionGroup{}, errors.New("failed to get inserted id")
	}

	mongoTransactionGroup.Id = insertedId

	return mongoTransactionGroup.ToDomain(), nil
}

// DeleteTransactionGroup implements domain.TransactionGroupRepository.
func (m MongoTransactionGroupRepository) DeleteTransactionGroup(transactionGroupId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(transactionGroupId)
	if err != nil {
		return err
	}

	return adapter.DeleteOne(client, ctx, transactionGroupCollection, bson.M{"_id": mongoId})
}

// GetTransactionGroup implements domain.TransactionGroupRepository.
func (m MongoTransactionGroupRepository) GetTransactionGroup(transactionGroupId string) (domain.TransactionGroup, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.TransactionGroup{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(transactionGroupId)
	if err != nil {
		return domain.TransactionGroup{}, err
	}

	filter := bson.M{"_id": mongoId}

	cursor, err := adapter.Query(client, ctx, transactionGroupCollection, filter, nil)
	if err != nil {
		return domain.TransactionGroup{}, err
	}

	mongoTransactionGroups := []MongoTransactionGroup{}
	err = cursor.All(ctx, &mongoTransactionGroups)
	if err != nil {
		return domain.TransactionGroup{}, err
	}

	if len(mongoTransactionGroups) == 0 {
		return domain.TransactionGroup{}, ErrorTransactionGroupNotFound
	}

	return mongoTransactionGroups[0].ToDomain(), nil
}

// GetTransactionGroupsForUser implements domain.TransactionGroupRepository.
func (m MongoTransactionGroupRepository) GetTransactionGroupsForUser(userId string) ([]domain.TransactionGroup, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": mongoId}

	cursor, err := adapter.Query(client, ctx, transactionGroupCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoTransactionGroups := []MongoTransactionGroup{}
	err = cursor.All(ctx, &mongoTransactionGroups)
	if err != nil {
		return nil, err
	}

	transactionGroups := make([]domain.TransactionGroup, len(mongoTransactionGroups))
	for i, mongoTransactionGroup := range mongoTransactionGroups {
		transactionGroups[i] = mongoTransactionGroup.ToDomain()
	}

	return transactionGroups, nil
}

// UpdateTransactionGroup implements domain.TransactionGroupRepository.
func (m MongoTransactionGroupRepository) UpdateTransactionGroup(transactionGroup domain.TransactionGroup) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoTransactionGroup, err := MongoTransactionGroupFromDomain(transactionGroup)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mongoTransactionGroup.Id}

	update := bson.M{"$set": mongoTransactionGroup}

	err = adapter.UpdateOne(client, ctx, transactionGroupCollection, filter, update)
	if err != nil {
		return err
	}

	return nil
}

var _ domain.TransactionGroupRepository = MongoTransactionGroupRepository{}
