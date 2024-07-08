package media_user_int

import (
	"encoding/json"
	"fmt"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/domain"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"io"
	"net/http"
	"os"
)

type UploadForm struct {
	OriginalFileName string `json:"originalFileName"`
	Extension        string `json:"extension"`
	FileBytes        []byte `json:"fileBytes"`
}

// returns new db file
func (c Controller) Upload(rw http.ResponseWriter, req *http.Request) {
	// decode body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var upload UploadForm
	err = json.Unmarshal(body, &upload)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// create new file in db
	currentUser, err := sessions_user_int.CurrentUser(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	newFile := domain.NewMediaFile(currentUser.Id, upload.Extension, upload.OriginalFileName)
	file, err := c.MediaRepository.Create(newFile)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// create new file in storage
	osFile, err := os.Create(fmt.Sprintf("%v/%v.%v", c.Config.Directory, file.Id, file.Extension))
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	defer osFile.Close()
	_, err = osFile.Write(upload.FileBytes)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// write new db file object to response writer
	res, err := json.Marshal(file)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}
