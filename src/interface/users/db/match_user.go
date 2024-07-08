package usersdb

import (
	"fmt"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/domain"
	"golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (r MongoUserRepository) MatchUser(emailOrUsername, password string) (domain.User, error) {
	// connect the db
	client, context, cancelFunc, err := mongo.Connect(r.connectionString())
	if err != nil {
		return domain.User{}, err
	}
	defer mongo.Close(client, context, cancelFunc)

	// find all email matches
	filter := bson.M{
		"email": emailOrUsername,
	}
	query := UserAgg(filter, bson.M{}, 0)
	emailMatches := []domain.User{}
	err = mongo.AggregatedQuery(client, context, r.DbName, r.collName(), query, &emailMatches)
	if err != nil {
		return domain.User{}, err
	}

	// find all username matches
	filter = bson.M{
		"username": emailOrUsername,
	}
	query = UserAgg(filter, bson.M{}, 0)
	usernameMatches := []domain.User{}
	err = mongo.AggregatedQuery(client, context, r.DbName, r.collName(), query, &usernameMatches)
	if err != nil {
		return domain.User{}, err
	}

	// check email matches & return first match
	for _, record := range emailMatches {
		if bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(password)) == nil {
			return record, nil
		}
	}

	// check username matches & return first match
	for _, record := range usernameMatches {
		if bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(password)) == nil {
			return record, nil
		}
	}

	// return error if there is no match
	return domain.User{}, fmt.Errorf(srverr.InvalidCredentials)
}
