package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type TransactionGroup struct {
	Id          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Description string        `json:"description" bson:"description"`
	UserId      bson.ObjectID `json:"userId" bson:"userId"`
}

func NewTransactionGroup(description string, userId bson.ObjectID) TransactionGroup {
	return TransactionGroup{
		Description: description,
		UserId:      userId,
	}
}
