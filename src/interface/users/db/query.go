package usersdb

import (
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r MongoUserRepository) QueryByEmail(search string) ([]domain.User, error) {
	// connect to the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return []domain.User{}, err
	}
	defer database.Close(client, context, cancelFunc)

	// get the records from the db
	filter := bson.M{
		"email": bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
			},
		},
	}
	query := UserAgg(filter, bson.M{}, 100)
	var records []domain.User
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Users, query, &records)
	if err != nil {
		return []domain.User{}, err
	}
	return records, nil
}

func (r MongoUserRepository) QueryByUsername(search string) ([]domain.User, error) {
	// connect to the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return []domain.User{}, err
	}
	defer database.Close(client, context, cancelFunc)

	// get the records from the db
	filter := bson.M{
		"username": bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
			},
		},
	}
	query := UserAgg(filter, bson.M{}, 100)
	var records []domain.User
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Users, query, &records)
	if err != nil {
		return []domain.User{}, err
	}
	return records, nil
}
