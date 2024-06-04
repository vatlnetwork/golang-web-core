package mediadb

import (
	"golang-web-core/src/domain"
)

type MediaFileRecord struct {
	Id               string `bson:"id"`
	OwnerId          string `bson:"ownerId"`
	Extension        string `bson:"extension"`
	OriginalFileName string `bson:"originalFileName"`
}

func MediaFileRecordFromDomain(file domain.MediaFile) MediaFileRecord {
	return MediaFileRecord{
		Id:               file.Id,
		OwnerId:          file.OwnerId,
		Extension:        file.Extension,
		OriginalFileName: file.OriginalFileName,
	}
}
