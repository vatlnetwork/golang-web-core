package mediadb

import (
	"fmt"
	"golang-web-core/src/domain"
	"golang-web-core/src/interface/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoMediaRepo struct {
	mongo.Config
}

func NewMongoMediaRepo(cfg mongo.Config) MongoMediaRepo {
	return MongoMediaRepo{
		Config: cfg,
	}
}

func (r MongoMediaRepo) collName() string {
	return r.CollNames.MediaFiles
}

func (r MongoMediaRepo) FindByUser(u string) ([]domain.MediaFile, error) {
	// connect to the db
	client, context, cancel, err := mongo.Connect(r.ConnectionString())
	if err != nil {
		return []domain.MediaFile{}, err
	}
	defer mongo.Close(client, context, cancel)

	// get the records
	query := bson.M{
		"ownerId": u,
	}
	var records []domain.MediaFile
	err = mongo.Query(client, context, r.DbName, r.collName(), query, &records, 0)
	if err != nil {
		return []domain.MediaFile{}, err
	}

	return records, nil
}

func (r MongoMediaRepo) Create(file domain.MediaFile) (domain.MediaFile, error) {
	// connect to the db
	client, context, cancel, err := mongo.Connect(r.ConnectionString())
	if err != nil {
		return domain.MediaFile{}, err
	}
	defer mongo.Close(client, context, cancel)

	// create the record
	record := MediaFileRecordFromDomain(file)
	err = mongo.InsertOne(client, context, r.DbName, r.collName(), record)
	if err != nil {
		return domain.MediaFile{}, err
	}

	// return the created record
	return file, nil
}

func (r MongoMediaRepo) Find(id string) (domain.MediaFile, error) {
	// connect to the db
	client, context, cancel, err := mongo.Connect(r.ConnectionString())
	if err != nil {
		return domain.MediaFile{}, err
	}
	defer mongo.Close(client, context, cancel)

	// find the record
	query := bson.M{
		"id": id,
	}
	records := []domain.MediaFile{}
	err = mongo.Query(client, context, r.DbName, r.collName(), query, &records, 1)
	if err != nil {
		return domain.MediaFile{}, err
	}
	if len(records) == 0 {
		return domain.MediaFile{}, fmt.Errorf("media record with id %v not found", id)
	}

	// return the record
	return records[0], nil
}

func (r MongoMediaRepo) Delete(id, u string) error {
	// find the record
	file, err := r.Find(id)
	if err != nil {
		return err
	}

	// verify ownership
	if file.OwnerId != u {
		return fmt.Errorf("no permission: you are not the owner of this file")
	}

	// connect the db
	client, context, cancel, err := mongo.Connect(r.ConnectionString())
	if err != nil {
		return err
	}
	defer mongo.Close(client, context, cancel)

	// delete the record
	record := MediaFileRecordFromDomain(file)
	filter := bson.M{
		"id": record.Id,
	}
	err = mongo.DeleteOne(client, context, r.DbName, r.collName(), filter)
	return err
}
