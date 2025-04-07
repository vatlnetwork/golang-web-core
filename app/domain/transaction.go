package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Transaction struct {
	Id          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount      float64       `json:"amount" bson:"amount"`
	Timestamp   int64         `json:"timestamp" bson:"timestamp"`
	Description string        `json:"description" bson:"description"`
	GroupId     string        `json:"groupId" bson:"groupId"`
}
