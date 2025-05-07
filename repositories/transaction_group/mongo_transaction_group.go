package transactiongrouprepo

import (
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoTransactionGroup struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	UserId      string        `bson:"userId"`
	Description string        `bson:"description"`
}

func (t MongoTransactionGroup) ToDomain() domain.TransactionGroup {
	return domain.TransactionGroup{
		Id:          t.Id.Hex(),
		UserId:      t.UserId,
		Description: t.Description,
	}
}

func MongoTransactionGroupFromDomain(transactionGroup domain.TransactionGroup) (MongoTransactionGroup, error) {
	mongoTransactionGroup := MongoTransactionGroup{
		UserId:      transactionGroup.UserId,
		Description: transactionGroup.Description,
	}

	if transactionGroup.Id != "" {
		id, err := bson.ObjectIDFromHex(transactionGroup.Id)
		if err != nil {
			return MongoTransactionGroup{}, err
		}
		mongoTransactionGroup.Id = id
	}

	return mongoTransactionGroup, nil
}
