package sessionsdb

import (
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoSessionRepository) QueryByUserId(userId string) ([]domain.Session, error) {
	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return []domain.Session{}, err
	}
	defer database.Close(client, context, cancelFunc)

	filter := bson.M{
		"userId": userId,
	}
	query := SessionAgg(filter, bson.M{}, 0)

	var records []domain.Session
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Sessions, query, &records)
	if err != nil {
		return []domain.Session{}, err
	}

	return records, nil
}
