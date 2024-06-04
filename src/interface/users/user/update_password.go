package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/domain"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/srverr"
	"io"
	"net/http"
)

type UpdatePasswordForm struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// returns a new session
func (c Controller) UpdatePassword(rw http.ResponseWriter, req *http.Request) {
	// decode request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var update UpdatePasswordForm
	err = json.Unmarshal(body, &update)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// get current session / user
	currentSessionPointer, err := sessions_user_int.CurrentSession(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	currentSession := *currentSessionPointer
	currentUser := currentSession.User

	// check current password
	_, err = c.UserRepository.MatchUser(currentUser.Email, update.CurrentPassword)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// update password
	currentUser.Password = update.NewPassword
	err = c.UserRepository.Update(currentUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// invalidate all sessions
	err = c.SessionRepository.DeleteAll(currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// create a new session
	newSession, err := domain.NewSession(currentUser, req.RemoteAddr, currentSession.DoesExpire)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	session, err := c.SessionRepository.FindOrCreate(newSession)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// return the new session
	res, err := json.Marshal(session)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}
