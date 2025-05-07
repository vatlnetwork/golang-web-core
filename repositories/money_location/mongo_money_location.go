package moneylocationrepo

import (
	"errors"
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoMoneyLocation struct {
	Id     bson.ObjectID `bson:"_id,omitempty"`
	UserId string        `bson:"userId"`
	Name   string        `bson:"name"`
}

func (m MongoMoneyLocation) ToDomain() domain.MoneyLocation {
	return domain.MoneyLocation{
		Id:     m.Id.Hex(),
		UserId: m.UserId,
		Name:   m.Name,
	}
}

func MongoMoneyLocationFromDomain(moneyLocation domain.MoneyLocation) (MongoMoneyLocation, error) {
	mongoMoneyLocation := MongoMoneyLocation{
		UserId: moneyLocation.UserId,
		Name:   moneyLocation.Name,
	}

	if moneyLocation.Id != "" {
		id, err := bson.ObjectIDFromHex(moneyLocation.Id)
		if err != nil {
			return MongoMoneyLocation{}, err
		}

		if id.IsZero() {
			return MongoMoneyLocation{}, errors.New("invalid money location id")
		}

		mongoMoneyLocation.Id = id
	}

	return mongoMoneyLocation, nil
}
