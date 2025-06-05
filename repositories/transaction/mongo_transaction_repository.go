package transactionrepo

import (
	"errors"
	"inventory-app/domain"
	"inventory-app/util/database_adapters/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const transactionCollection string = "transactions"

type MongoTransactionRepository struct {
	connectionConfig mongo.Config
	logTransactions  bool
}

func NewMongoTransactionRepository(connectionConfig mongo.Config, logTransactions bool) MongoTransactionRepository {
	return MongoTransactionRepository{connectionConfig: connectionConfig, logTransactions: logTransactions}
}

func (m MongoTransactionRepository) adapter() *mongo.Mongo {
	return mongo.NewMongoAdapter(m.connectionConfig, m.logTransactions)
}

// DeleteAllTransactionsForUser implements domain.TransactionRepository.
func (m MongoTransactionRepository) DeleteAllTransactionsForUser(userId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	err = adapter.DeleteMany(client, ctx, transactionCollection, filter)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTransactionsInGroup implements domain.TransactionRepository.
func (m MongoTransactionRepository) DeleteTransactionsInGroup(groupId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"groupId": groupId}

	return adapter.DeleteMany(client, ctx, transactionCollection, filter)
}

// DeleteTransactionsInLocation implements domain.TransactionRepository.
func (m MongoTransactionRepository) DeleteTransactionsInLocation(locationId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"moneyLocationId": locationId}

	return adapter.DeleteMany(client, ctx, transactionCollection, filter)
}

// GetTransactionsByGroup implements domain.TransactionRepository.
func (m MongoTransactionRepository) GetTransactionsByGroup(groupId string) ([]domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"groupId": groupId}

	cursor, err := adapter.Query(client, ctx, transactionCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoTransactions := []MongoTransaction{}
	err = cursor.All(ctx, &mongoTransactions)
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, len(mongoTransactions))
	for i, mongoTransaction := range mongoTransactions {
		transactions[i] = mongoTransaction.ToDomain()
	}

	return transactions, nil
}

// GetTransactionsByLocation implements domain.TransactionRepository.
func (m MongoTransactionRepository) GetTransactionsByLocation(locationId string) ([]domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"moneyLocationId": locationId}

	cursor, err := adapter.Query(client, ctx, transactionCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoTransactions := []MongoTransaction{}
	err = cursor.All(ctx, &mongoTransactions)
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, len(mongoTransactions))
	for i, mongoTransaction := range mongoTransactions {
		transactions[i] = mongoTransaction.ToDomain()
	}

	return transactions, nil
}

// GetTransactionsByYear implements domain.TransactionRepository.
func (m MongoTransactionRepository) GetTransactionsByYear(userId string, year int) ([]domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId, "year": year}

	cursor, err := adapter.Query(client, ctx, transactionCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoTransactions := []MongoTransaction{}
	err = cursor.All(ctx, &mongoTransactions)
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, len(mongoTransactions))
	for i, mongoTransaction := range mongoTransactions {
		transactions[i] = mongoTransaction.ToDomain()
	}

	return transactions, nil
}

// CreateTransaction implements domain.TransactionRepository.
func (m MongoTransactionRepository) CreateTransaction(transaction domain.Transaction) (domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.Transaction{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoTransaction, err := MongoTransactionFromDomain(transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	result, err := adapter.InsertOne(client, ctx, transactionCollection, mongoTransaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	insertedId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.Transaction{}, errors.New("failed to get inserted id")
	}

	mongoTransaction.Id = insertedId

	return mongoTransaction.ToDomain(), nil
}

// DeleteTransaction implements domain.TransactionRepository.
func (m MongoTransactionRepository) DeleteTransaction(transactionId string) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(transactionId)
	if err != nil {
		return err
	}

	return adapter.DeleteOne(client, ctx, transactionCollection, bson.M{"_id": mongoId})
}

// GetTransaction implements domain.TransactionRepository.
func (m MongoTransactionRepository) GetTransaction(transactionId string) (domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return domain.Transaction{}, err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoId, err := bson.ObjectIDFromHex(transactionId)
	if err != nil {
		return domain.Transaction{}, err
	}

	filter := bson.M{"_id": mongoId}

	cursor, err := adapter.Query(client, ctx, transactionCollection, filter, nil)
	if err != nil {
		return domain.Transaction{}, err
	}

	mongoTransactions := []MongoTransaction{}
	err = cursor.All(ctx, &mongoTransactions)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(mongoTransactions) == 0 {
		return domain.Transaction{}, errors.New(domain.ErrorTransactionNotFound)
	}

	return mongoTransactions[0].ToDomain(), nil
}

// GetTransactionsForUser implements domain.TransactionRepository.
func (m MongoTransactionRepository) GetTransactionsForUser(userId string) ([]domain.Transaction, error) {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	defer adapter.Close(client, ctx, cancel)

	filter := bson.M{"userId": userId}

	cursor, err := adapter.Query(client, ctx, transactionCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	mongoTransactions := []MongoTransaction{}
	err = cursor.All(ctx, &mongoTransactions)
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, len(mongoTransactions))
	for i, mongoTransaction := range mongoTransactions {
		transactions[i] = mongoTransaction.ToDomain()
	}

	return transactions, nil
}

// UpdateTransaction implements domain.TransactionRepository.
func (m MongoTransactionRepository) UpdateTransaction(transaction domain.Transaction) error {
	adapter := m.adapter()

	client, ctx, cancel, err := adapter.Connect()
	if err != nil {
		return err
	}
	defer adapter.Close(client, ctx, cancel)

	mongoTransaction, err := MongoTransactionFromDomain(transaction)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mongoTransaction.Id}

	update := bson.M{"$set": mongoTransaction}

	err = adapter.UpdateOne(client, ctx, transactionCollection, filter, update)
	if err != nil {
		return err
	}

	return nil
}

var _ domain.TransactionRepository = MongoTransactionRepository{}
