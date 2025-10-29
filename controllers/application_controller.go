package controllers

import (
	"errors"
	"golang-web-core/domain"
	"golang-web-core/services/httpserver"
	"net/http"
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
