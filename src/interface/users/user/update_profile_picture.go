package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"io"
	"net/http"
)

type UpdateProfilePictureForm struct {
	MediaFileId string `json:"mediaFileId"`
}

func (c Controller) UpdateProfilePicture(rw http.ResponseWriter, req *http.Request) {
	// decode the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var update UpdateProfilePictureForm
	err = json.Unmarshal(body, &update)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// get current user
	currentUser, err := sessions_user_int.CurrentUser(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// find the image
	mediaFile, err := c.MediaRepository.Find(update.MediaFileId)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// update the current user's profile picture
	currentUser.Image = mediaFile
	err = c.UserRepository.Update(currentUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
