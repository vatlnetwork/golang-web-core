package sessionsdb

import (
	database "golang-web-core/src/interface/database/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoSessionRepository) Delete(id string) error {
	filter := bson.M{
		"id": id,
	}

	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancelFunc)
	err = database.DeleteOne(client, context, r.Config.DbName, r.Config.CollNames.Sessions, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r MongoSessionRepository) DeleteExpired(userId string) error {
	filter := bson.M{
		"userId":     userId,
		"doesExpire": true,
		"expiresAt": bson.M{
			"$lt": time.Now().UnixMilli(),
		},
	}

	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancelFunc)
	err = database.DeleteMany(client, context, r.Config.DbName, r.Config.CollNames.Sessions, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r MongoSessionRepository) DeleteAll(userId string) error {
	// connect the db
	client, context, cancel, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancel)

	// delete the records
	filter := bson.M{
		"userId": userId,
	}
	err = database.DeleteMany(client, context, r.DbName, r.collName(), filter)
	return err
}
