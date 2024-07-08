package sessions_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	"golang-web-core/src/application/srv/routes"
	"golang-web-core/src/domain"
	"net/http"
)

type Controller struct {
	domain.SessionRepository
}

func NewSessionController(sessionRepo domain.SessionRepository) Controller {
	return Controller{
		SessionRepository: sessionRepo,
	}
}

func (c Controller) Routes() []routes.Route {
	return []routes.Route{
		{
			Path:    "/sessions/current",
			Method:  http.MethodGet,
			Handler: c.CurrentSession,
		},
		{
			Path:    "/sessions/delete",
			Method:  http.MethodDelete,
			Handler: c.DeleteSession,
		},
	}
}

func (c Controller) CurrentSession(rw http.ResponseWriter, req *http.Request) {
	session, err := CurrentSession(req)
	if session == nil || err != nil {
		rw.Write([]byte{})
		return
	}
	res, err := json.Marshal(*session)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}

func (c Controller) DeleteSession(rw http.ResponseWriter, req *http.Request) {
	session, err := CurrentSession(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	if session != nil {
		err := c.SessionRepository.DeleteExpired((*session).User.Id)
		if err != nil {
			srverr.Raise(rw, req, err, http.StatusInternalServerError)
			return
		}
		err = c.SessionRepository.Delete((*session).Id)
		if err != nil {
			srverr.Raise(rw, req, err, http.StatusInternalServerError)
			return
		}
	}
}
