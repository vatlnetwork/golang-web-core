package userrepo

import (
	"errors"
	"inventory-app/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoUser struct {
	Id                bson.ObjectID `bson:"_id,omitempty"`
	Email             string        `bson:"email"`
	EncryptedPassword string        `bson:"encryptedPassword"`
}

func (m MongoUser) ToDomain() domain.User {
	return domain.User{
		Id:                m.Id.Hex(),
		Email:             m.Email,
		EncryptedPassword: m.EncryptedPassword,
	}
}

func MongoUserFromDomain(user domain.User) (MongoUser, error) {
	mongoUser := MongoUser{
		Email:             user.Email,
		EncryptedPassword: user.EncryptedPassword,
	}

	if user.Id != "" {
		id, err := bson.ObjectIDFromHex(user.Id)
		if err != nil {
			return MongoUser{}, err
		}

		if id.IsZero() {
			return MongoUser{}, errors.New("invalid user id")
		}

		mongoUser.Id = id
	}

	return mongoUser, nil
}
