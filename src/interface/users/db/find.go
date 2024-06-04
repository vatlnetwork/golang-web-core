package usersdb

import (
	"fmt"
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoUserRepository) Find(id string) (domain.User, error) {
	// connect to the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return domain.User{}, err
	}
	defer database.Close(client, context, cancelFunc)

	// find record in the database
	query := UserAgg(bson.M{"id": id}, bson.M{}, 1)
	var records []domain.User
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Users, query, &records)
	if err != nil {
		return domain.User{}, err
	}

	// return record if it exists
	if len(records) == 0 {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return records[0], nil
}

func (r MongoUserRepository) FindByEmail(email string) (domain.User, error) {
	// connect to the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return domain.User{}, err
	}
	defer database.Close(client, context, cancelFunc)

	// find record in the database
	filter := bson.M{
		"email": email,
	}
	query := UserAgg(filter, bson.M{}, 1)
	var records []domain.User
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Users, query, &records)
	if err != nil {
		return domain.User{}, err
	}

	// return record if it exists
	if len(records) == 0 {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return records[0], nil
}
