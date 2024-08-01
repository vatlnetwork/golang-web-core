package controllers

import (
	"golang-web-core/srv/cfg"
	"net/http"
)

type ApplicationController struct {
	cfg.Config
}

func NewApplicationController(config cfg.Config) ApplicationController {
	return ApplicationController{
		Config: config,
	}
}

func (c ApplicationController) Name() string {
	return "ApplicationController"
}

func (c ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}
