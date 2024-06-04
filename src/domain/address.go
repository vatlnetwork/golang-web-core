package domain

type Address struct {
	Street  string `json:"street" bson:"street"`
	AptNo   string `json:"aptNo" bson:"aptNo"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Zip     int64  `json:"zip" bson:"zip"`
	Country string `json:"country" bson:"country"`
}
