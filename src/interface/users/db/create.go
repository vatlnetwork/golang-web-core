package usersdb

import (
	"fmt"
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoUserRepository) CreateUser(user domain.User) (domain.User, error) {
	// generate the record
	userRecord, err := UserRecordFromDomain(user)
	if err != nil {
		return domain.User{}, err
	}

	// connect the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return domain.User{}, err
	}
	defer database.Close(client, context, cancelFunc)

	// check to see if a record with the same email already exists
	filter := bson.M{
		"email": user.Email,
	}
	var records []UserRecord
	err = database.Query(client, context, r.Config.DbName, r.collName(), filter, &records, 1)
	if err != nil {
		return domain.User{}, err
	}
	if len(records) > 0 {
		return domain.User{}, fmt.Errorf("the provided email is already taken")
	}

	// check to see if there is an admin on the system and promote the user
	// if there is no admin on the system
	filter = bson.M{
		"isAdmin": true,
	}
	records = []UserRecord{}
	err = database.Query(client, context, r.Config.DbName, r.Config.CollNames.Users, filter, &records, 1)
	if err != nil {
		return domain.User{}, err
	}
	if len(records) == 0 {
		userRecord.IsAdmin = true
	}

	// insert & return the record
	err = database.InsertOne(client, context, r.Config.DbName, r.Config.CollNames.Users, userRecord)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
