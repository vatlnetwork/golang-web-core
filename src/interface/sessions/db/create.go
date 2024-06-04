package sessionsdb

import (
	"golang-web-core/src/domain"
	database "golang-web-core/src/interface/database/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (r MongoSessionRepository) FindOrCreate(session domain.Session) (domain.Session, error) {
	sessionRecord := SessionRecordFromDomain(session)

	client, context, cancelFunc, err := database.Connect(r.connectionString())
	if err != nil {
		return domain.Session{}, err
	}
	defer database.Close(client, context, cancelFunc)

	filter := bson.M{
		"userId":     session.User.Id,
		"verifiedIp": session.VerifiedIp,
		"expiresAt": bson.M{
			"$gt": time.Now().UnixMilli(),
		},
	}
	query := SessionAgg(filter, bson.M{}, 0)
	var records []domain.Session
	err = database.AggregatedQuery(client, context, r.Config.DbName, r.Config.CollNames.Sessions, query, &records)
	if err != nil {
		return domain.Session{}, err
	}
	for _, record := range records {
		if !record.IsExpired() {
			return record, nil
		}
	}

	err = database.InsertOne(client, context, r.Config.DbName, r.Config.CollNames.Sessions, sessionRecord)
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}
