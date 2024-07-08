package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/domain"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"io"
	"net/http"
)

func (c Controller) UpdateAddress(rw http.ResponseWriter, req *http.Request) {
	// decode request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// get new address out of request body
	var address domain.Address
	err = json.Unmarshal(body, &address)
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

	// update user address
	currentUser.Address = address
	err = c.UserRepository.Update(currentUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
