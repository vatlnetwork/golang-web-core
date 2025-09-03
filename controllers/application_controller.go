package controllers

import (
	"golang-web-core/logging"
	"golang-web-core/services/httpserver"
	"net/http"
)

type ApplicationController struct {
	logger *logging.Logger
}

func NewApplicationController(logger *logging.Logger) ApplicationController {
	return ApplicationController{
		logger: logger,
	}
}

// BeforeAction implements Controller.
func (a ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

// Routes implements Controller.
func (a ApplicationController) Routes() []httpserver.Route {
	return []httpserver.Route{}
}

var _ Controller = ApplicationController{}
