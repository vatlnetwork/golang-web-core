package transactionrepo

import (
	"errors"
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoTransaction struct {
	Id              bson.ObjectID `bson:"_id,omitempty"`
	UserId          string        `bson:"userId"`
	Amount          float64       `bson:"amount"`
	Timestamp       int64         `bson:"timestamp"`
	Year            int           `bson:"year"`
	Description     string        `bson:"description,omitempty"`
	GroupId         string        `bson:"groupId,omitempty"`
	MoneyLocationId string        `bson:"moneyLocationId"`
}

func (t MongoTransaction) ToDomain() domain.Transaction {
	return domain.Transaction{
		Id:              t.Id.Hex(),
		UserId:          t.UserId,
		Amount:          t.Amount,
		Timestamp:       t.Timestamp,
		Year:            t.Year,
		Description:     t.Description,
		GroupId:         t.GroupId,
		MoneyLocationId: t.MoneyLocationId,
	}
}

func MongoTransactionFromDomain(transaction domain.Transaction) (MongoTransaction, error) {
	mongoTransaction := MongoTransaction{
		UserId:          transaction.UserId,
		Amount:          transaction.Amount,
		Timestamp:       transaction.Timestamp,
		Year:            transaction.Year,
		Description:     transaction.Description,
		GroupId:         transaction.GroupId,
		MoneyLocationId: transaction.MoneyLocationId,
	}

	if transaction.Id != "" {
		id, err := bson.ObjectIDFromHex(transaction.Id)
		if err != nil {
			return MongoTransaction{}, err
		}

		if id.IsZero() {
			return MongoTransaction{}, errors.New("invalid transaction id")
		}

		mongoTransaction.Id = id
	}

	return mongoTransaction, nil
}
