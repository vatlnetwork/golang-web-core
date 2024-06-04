package usersdb

import (
	"fmt"
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoUserRepository) Update(user domain.User) error {
	// convert user to user record
	userRecord, err := UserRecordFromDomain(user)
	if err != nil {
		return err
	}

	// connect the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancelFunc)

	// find current user
	currentUser, err := r.Find(userRecord.Id)
	if err != nil {
		return err
	}

	// make sure new email doesn't already exist
	if userRecord.Email != currentUser.Email {
		emailFilter := bson.M{
			"email": userRecord.Email,
		}
		var results []UserRecord
		err = database.Query(client, context, r.DbName, r.collName(), emailFilter, &results, 0)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			return fmt.Errorf("a user with that email already exists")
		}
	}

	// update record
	filter := bson.M{
		"id": userRecord.Id,
	}
	update := bson.M{
		"$set": userRecord,
	}
	err = database.UpdateOne(client, context, r.Config.DbName, r.Config.CollNames.Users, filter, update)
	if err != nil {
		return err
	}
	return nil
}
