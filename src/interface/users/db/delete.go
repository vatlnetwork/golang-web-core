package usersdb

import (
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoUserRepository) Delete(id string) error {
	// connect to the db
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancelFunc)

	// delete the record
	filter := bson.M{
		"id": id,
	}
	err = database.DeleteOne(client, context, r.Config.DbName, r.Config.CollNames.Users, filter)
	if err != nil {
		return err
	}
	return nil
}
