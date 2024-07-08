package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/domain"
	"io"
	"net/http"
	"time"
)

type SignInForm struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
	ExpireSession   bool   `json:"expireSession"`
}

// returns a new session
func (c Controller) SignIn(rw http.ResponseWriter, req *http.Request) {
	// decode req body
	var signInDetails SignInForm
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &signInDetails)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// match user
	user, err := c.UserRepository.MatchUser(signInDetails.EmailOrUsername, signInDetails.Password)
	if err != nil {
		if err.Error() == srverr.InvalidCredentials {
			srverr.Raise(rw, req, err, http.StatusNotFound)
			return
		}
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// find or create new session
	newSession, err := domain.NewSession(user, req.RemoteAddr, signInDetails.ExpireSession)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	session, err := c.SessionRepository.FindOrCreate(newSession)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// update last sign in
	user.LastSignIn = time.Now().UnixMilli()
	err = c.UserRepository.Update(user)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// write the new session to the response body
	res, err := json.Marshal(session)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}
