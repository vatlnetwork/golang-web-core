package sessionsdb

import (
	"golang-web-core/src/interface/database/mongo"
)

type MongoSessionRepository struct {
	mongo.Config
}

func NewMongoSessionRepository(cfg mongo.Config) MongoSessionRepository {
	return MongoSessionRepository{
		Config: cfg,
	}
}

func (r MongoSessionRepository) connectionString() string {
	return r.Config.ConnectionString()
}

func (r MongoSessionRepository) collName() string {
	return r.CollNames.Sessions
}
