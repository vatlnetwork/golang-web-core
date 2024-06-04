package users_user_int

import (
	"encoding/json"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/srverr"
	"io"
	"net/http"
)

type ConfirmPasswordForm struct {
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (c Controller) ConfirmPassword(rw http.ResponseWriter, req *http.Request) {
	// get password from request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var form ConfirmPasswordForm
	err = json.Unmarshal(body, &form)
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

	// check password
	_, err = c.UserRepository.MatchUser(currentUser.Email, form.PasswordConfirmation)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
