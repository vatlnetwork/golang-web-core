package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	Id          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount      float64       `json:"amount" bson:"amount"`
	Timestamp   int64         `json:"timestamp" bson:"timestamp"`
	Description string        `json:"description" bson:"description"`
	GroupId     bson.ObjectID `json:"groupId" bson:"groupId,omitempty"`
	Year        int           `json:"year" bson:"year"`
}

func NewTransaction(amount float64, description string, groupId string) (Transaction, error) {
	var transactionGroupId bson.ObjectID
	if groupId != "" {
		var err error
		transactionGroupId, err = bson.ObjectIDFromHex(groupId)
		if err != nil {
			return Transaction{}, err
		}
	}

	transaction := Transaction{
		Amount:      amount,
		Timestamp:   time.Now().UnixMilli(),
		Description: description,
		GroupId:     transactionGroupId,
		Year:        time.Now().Year(),
	}

	return transaction, nil
}
