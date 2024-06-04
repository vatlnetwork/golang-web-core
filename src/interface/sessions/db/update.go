package sessionsdb

import (
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoSessionRepository) Update(session domain.Session) error {
	sessionRecord := SessionRecordFromDomain(session)
	filter := bson.M{
		"id": sessionRecord.Id,
	}
	update := bson.M{
		"$set": sessionRecord,
	}
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return err
	}
	defer database.Close(client, context, cancelFunc)
	err = database.UpdateOne(client, context, r.Config.DbName, r.Config.CollNames.Sessions, filter, update)
	if err != nil {
		return err
	}
	return nil
}
