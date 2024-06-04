package domain

import "github.com/google/uuid"

type MediaFile struct {
	Id               string `json:"id" bson:"id"`
	OwnerId          string `json:"ownerId" bson:"ownerId"`
	Extension        string `json:"extension" bson:"extension"`
	OriginalFileName string `json:"originalFileName" bson:"originalFileName"`
}

func NewMediaFile(owner, extension, originalFileName string) MediaFile {
	return MediaFile{
		Id:               uuid.NewString(),
		OwnerId:          owner,
		Extension:        extension,
		OriginalFileName: originalFileName,
	}
}
