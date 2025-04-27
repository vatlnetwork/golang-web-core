package transactiongrouprepo

import (
	"errors"
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoTransactionGroup struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	UserId      bson.ObjectID `bson:"userId"`
	Description string        `bson:"description"`
}

func (t MongoTransactionGroup) ToDomain() domain.TransactionGroup {
	return domain.TransactionGroup{
		Id:          t.Id.Hex(),
		UserId:      t.UserId.Hex(),
		Description: t.Description,
	}
}

func MongoTransactionGroupFromDomain(transactionGroup domain.TransactionGroup) (MongoTransactionGroup, error) {
	mongoTransactionGroup := MongoTransactionGroup{
		Description: transactionGroup.Description,
	}

	if transactionGroup.Id != "" {
		id, err := bson.ObjectIDFromHex(transactionGroup.Id)
		if err != nil {
			return MongoTransactionGroup{}, err
		}
		mongoTransactionGroup.Id = id
	}

	if transactionGroup.UserId != "" {
		id, err := bson.ObjectIDFromHex(transactionGroup.UserId)
		if err != nil {
			return MongoTransactionGroup{}, err
		}
		if id.IsZero() {
			return MongoTransactionGroup{}, errors.New("invalid user id")
		}
		mongoTransactionGroup.UserId = id
	}

	return mongoTransactionGroup, nil
}
