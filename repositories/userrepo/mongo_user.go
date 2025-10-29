package userrepo

import (
	"errors"
	"golang-web-core/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoUser struct {
	Id                bson.ObjectID `bson:"_id,omitempty"`
	Email             string        `bson:"email"`
	FirstName         string        `bson:"firstName"`
	LastName          string        `bson:"lastName"`
	EncryptedPassword string        `bson:"encryptedPassword"`
	CreatedAt         time.Time     `bson:"createdAt"`
	UpdatedAt         time.Time     `bson:"updatedAt"`
	LastSignIn        time.Time     `bson:"lastSignIn"`
}

func (m MongoUser) ToDomain() domain.User {
	return domain.User{
		Id:                m.Id.Hex(),
		Email:             m.Email,
		FirstName:         m.FirstName,
		LastName:          m.LastName,
		EncryptedPassword: m.EncryptedPassword,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		LastSignIn:        m.LastSignIn,
	}
}

func MongoUserFromDomain(user domain.User) (MongoUser, error) {
	mongoUser := MongoUser{
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		EncryptedPassword: user.EncryptedPassword,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		LastSignIn:        user.LastSignIn,
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
