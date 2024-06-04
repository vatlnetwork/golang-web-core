package domain

type Color struct {
	R int64 `json:"r" bson:"r"`
	G int64 `json:"g" bson:"g"`
	B int64 `json:"b" bson:"b"`
}
