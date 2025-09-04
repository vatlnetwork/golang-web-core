package controllers

import (
	"errors"
	"golang-web-core/logging"
	"golang-web-core/services/httpserver"
	"net/http"
)

type ApplicationController struct {
	logger *logging.Logger
}

func NewApplicationController(logger *logging.Logger) (ApplicationController, error) {
	if logger == nil {
		return ApplicationController{}, errors.New("logger is required")
	}

	return ApplicationController{
		logger: logger,
	}, nil
}

// BeforeAction implements httpserver.Controller.
func (a ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

// Routes implements httpserver.Controller.
func (a ApplicationController) Routes() []httpserver.Route {
	return []httpserver.Route{}
}

var _ httpserver.Controller = ApplicationController{}
