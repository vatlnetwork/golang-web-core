package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	Id          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount      float64       `json:"amount" bson:"amount"`
	Timestamp   time.Time     `json:"timestamp" bson:"timestamp"`
	Description string        `json:"description" bson:"description"`
	GroupId     bson.ObjectID `json:"groupId" bson:"groupId,omitempty"`
	Year        int           `json:"year" bson:"year"`
	UserId      bson.ObjectID `json:"userId" bson:"userId"`
}

func NewTransaction(amount float64, description, groupId string, userId bson.ObjectID) (Transaction, error) {
	var transactionGroupId bson.ObjectID
	if groupId != "" {
		var err error
		transactionGroupId, err = bson.ObjectIDFromHex(groupId)
		if err != nil {
			return Transaction{}, err
		}
	}

	if userId.IsZero() {
		return Transaction{}, fmt.Errorf("user id is required")
	}

	transaction := Transaction{
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: description,
		GroupId:     transactionGroupId,
		Year:        time.Now().Year(),
		UserId:      userId,
	}

	return transaction, nil
}
