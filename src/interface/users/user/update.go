package users_user_int

import (
	"encoding/json"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/srverr"
	"io"
	"net/http"
)

type UpdateForm struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	PhoneNo  int64  `json:"phoneNo"`
	PayCode  string `json:"payCode"`
}

func (c Controller) Update(rw http.ResponseWriter, req *http.Request) {
	// decode request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var update UpdateForm
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

	// update current user
	currentUser.Email = update.Email
	currentUser.Username = update.Username
	currentUser.PhoneNo = update.PhoneNo
	err = c.UserRepository.Update(currentUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
