package transactionrepo

import (
	"errors"
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoTransaction struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	UserId      bson.ObjectID `bson:"userId"`
	Amount      float64       `bson:"amount"`
	Timestamp   int64         `bson:"timestamp"`
	Year        int           `bson:"year"`
	Description string        `bson:"description,omitempty"`
	GroupId     bson.ObjectID `bson:"groupId,omitempty"`
}

func (t MongoTransaction) ToDomain() domain.Transaction {
	return domain.Transaction{
		Id:          t.Id.Hex(),
		UserId:      t.UserId.Hex(),
		Amount:      t.Amount,
		Timestamp:   t.Timestamp,
		Year:        t.Year,
		Description: t.Description,
		GroupId:     t.GroupId.Hex(),
	}
}

func MongoTransactionFromDomain(transaction domain.Transaction) (MongoTransaction, error) {
	mongoTransaction := MongoTransaction{
		Amount:      transaction.Amount,
		Timestamp:   transaction.Timestamp,
		Year:        transaction.Year,
		Description: transaction.Description,
	}

	if transaction.Id != "" {
		id, err := bson.ObjectIDFromHex(transaction.Id)
		if err != nil {
			return MongoTransaction{}, err
		}
		mongoTransaction.Id = id
	}

	if transaction.UserId != "" {
		id, err := bson.ObjectIDFromHex(transaction.UserId)
		if err != nil {
			return MongoTransaction{}, err
		}
		if id.IsZero() {
			return MongoTransaction{}, errors.New("invalid user id")
		}
		mongoTransaction.UserId = id
	}

	if transaction.GroupId != "" {
		id, err := bson.ObjectIDFromHex(transaction.GroupId)
		if err != nil {
			return MongoTransaction{}, err
		}
		mongoTransaction.GroupId = id
	}

	return mongoTransaction, nil
}
