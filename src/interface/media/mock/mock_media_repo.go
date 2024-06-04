package media_mock_int

import (
	"fmt"
	"golang-web-core/src/domain"

	"github.com/google/uuid"
)

type MockMediaRepo struct{}

func NewMockMediaRepo() MockMediaRepo {
	return MockMediaRepo{}
}

var mediaFiles []domain.MediaFile = []domain.MediaFile{}

func (r MockMediaRepo) FindByUser(u string) ([]domain.MediaFile, error) {
	res := []domain.MediaFile{}
	for _, file := range mediaFiles {
		if file.OwnerId == u {
			res = append(res, file)
		}
	}
	return res, nil
}

func (r MockMediaRepo) Create(file domain.MediaFile) (domain.MediaFile, error) {
	file.Id = uuid.NewString()
	mediaFiles = append(mediaFiles, file)
	return file, nil
}

func (r MockMediaRepo) Find(id string) (domain.MediaFile, error) {
	for _, file := range mediaFiles {
		if file.Id == id {
			return file, nil
		}
	}
	return domain.MediaFile{}, fmt.Errorf("unable to find media file with id %v", id)
}

func (r MockMediaRepo) Delete(id, u string) error {
	newFiles := []domain.MediaFile{}
	user := ""
	for _, file := range mediaFiles {
		if file.Id != id {
			newFiles = append(newFiles, file)
		} else {
			user = file.OwnerId
		}
	}
	if u != user {
		return fmt.Errorf("no permission")
	}
	mediaFiles = newFiles
	return nil
}
