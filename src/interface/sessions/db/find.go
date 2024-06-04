package sessionsdb

import (
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r MongoSessionRepository) Find(id string) (domain.Session, error) {
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return domain.Session{}, err
	}
	defer database.Close(client, context, cancelFunc)

	filter := bson.M{"id": id}
	query := SessionAgg(filter, bson.M{}, 1)

	var records []domain.Session
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Sessions, query, &records)
	if err != nil {
		return domain.Session{}, err
	}

	if len(records) == 0 {
		return domain.Session{}, mongo.ErrNoDocuments
	}
	record := records[0]
	return record, nil
}
