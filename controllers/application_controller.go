package controllers

import (
	"errors"
	"golang-web-core/domain"
	"golang-web-core/services/httpserver"
	"net/http"
	"strings"
)

type ApplicationController struct {
	errorHandler   *httpserver.HttpErrorHandler
	sessionManager domain.SessionManager
}

func NewApplicationController(errorHandler *httpserver.HttpErrorHandler, sessionManager domain.SessionManager) (ApplicationController, error) {
	if errorHandler == nil {
		return ApplicationController{}, errors.New("error handler is required")
	}

	return ApplicationController{
		errorHandler:   errorHandler,
		sessionManager: sessionManager,
	}, nil
}

// BeforeAction implements httpserver.Controller.
func (a ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// this skips checking auth for files to prevent tons of requests to the database
		// if you have an api endpoint with a . here that needs auth we will have to manually handle it
		if strings.Contains(r.URL.Path, ".") {
			handler(w, r)
			return
		}

		currentSession, currentUser, err := a.sessionManager.GetCurrentSession(r)
		if err != nil {
			a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
			return
		}
		if currentSession != nil {
			r = a.sessionManager.SetContextUser(r, currentUser)
		}

		handler(w, r)
	}
}

// Routes implements httpserver.Controller.
func (a ApplicationController) Routes() []httpserver.Route {
	return []httpserver.Route{}
}

var _ httpserver.Controller = ApplicationController{}
