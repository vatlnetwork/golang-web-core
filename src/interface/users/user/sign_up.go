package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/domain"
	"io"
	"net/http"
)

type SignUpForm struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ExpireSession bool   `json:"expireSession"`
}

// returns a new session
func (c Controller) SignUp(rw http.ResponseWriter, req *http.Request) {
	// decode body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	var signUpDetails SignUpForm
	err = json.Unmarshal(body, &signUpDetails)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// create new user
	newUser := domain.NewUser(signUpDetails.Email, signUpDetails.Username, signUpDetails.Password, false)
	user, err := c.UserRepository.CreateUser(newUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// create new session
	newSession, err := domain.NewSession(user, req.RemoteAddr, signUpDetails.ExpireSession)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	session, err := c.SessionRepository.FindOrCreate(newSession)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// write new session to the body
	res, err := json.Marshal(session)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}
