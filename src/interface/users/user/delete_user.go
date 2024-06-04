package users_user_int

import (
	"encoding/json"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/srverr"
	"io"
	"net/http"
)

type DeleteUserForm struct {
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (c Controller) DeleteUser(rw http.ResponseWriter, req *http.Request) {
	// get password confirmation out of request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var form DeleteUserForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	passwordConfirm := form.PasswordConfirmation

	// get current user
	currentUser, err := sessions_user_int.CurrentUser(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// match the user with the password
	_, err = c.UserRepository.MatchUser(currentUser.Email, passwordConfirm)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// delete all of the user's media
	media, err := c.MediaRepository.FindByUser(currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	for _, file := range media {
		err = c.MediaRepository.Delete(file.Id, currentUser.Id)
		if err != nil {
			srverr.Raise(rw, req, err, http.StatusInternalServerError)
			return
		}
		err = c.OsMediaRepo.Delete(file, currentUser.Id)
		if err != nil {
			srverr.Raise(rw, req, err, http.StatusInternalServerError)
			return
		}
	}

	// delete all of the user's sessions
	err = c.SessionRepository.DeleteAll(currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// delete the user
	err = c.UserRepository.Delete(currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
