package usersdb

import (
	"golang-web-core/src/interface/database/mongo"
)

type MongoUserRepository struct {
	mongo.Config
}

func NewMongoUserRepository(cfg mongo.Config) MongoUserRepository {
	return MongoUserRepository{
		Config: cfg,
	}
}

func (r MongoUserRepository) connectionString() string {
	return r.Config.ConnectionString()
}

func (r MongoUserRepository) collName() string {
	return r.CollNames.Users
}
