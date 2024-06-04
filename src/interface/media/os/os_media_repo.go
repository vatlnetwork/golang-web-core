package media_os_int

import (
	"fmt"
	"golang-web-core/src/domain"
	"os"
)

type OsMediaRepo struct {
	Config
}

func NewOsMediaRepo(cfg Config) OsMediaRepo {
	return OsMediaRepo{
		Config: cfg,
	}
}

func (r OsMediaRepo) Create(file domain.MediaFile, data []byte) error {
	// create file with id and extension at media path
	osFile, err := os.Create(fmt.Sprintf("%v/%v.%v", r.Directory, file.Id, file.Extension))
	if err != nil {
		return err
	}
	defer osFile.Close()
	_, err = osFile.Write(data)
	return err
}

func (r OsMediaRepo) Load(file domain.MediaFile) ([]byte, error) {
	// load file with id and extension from media path
	return os.ReadFile(fmt.Sprintf("%v/%v.%v", r.Directory, file.Id, file.Extension))
}

func (r OsMediaRepo) Delete(file domain.MediaFile, u string) error {
	// verify ownership
	if file.OwnerId != u {
		return fmt.Errorf("no permission: you do not own this file")
	}

	// delete file with id and extension from media path
	return os.Remove(fmt.Sprintf("%v/%v.%v", r.Directory, file.Id, file.Extension))
}
