package controllers

import (
	"golang-web-core/srv"
	"net/http"
)

type ApplicationController struct {
	srv.Config
}

func NewApplicationController(config srv.Config) ApplicationController {
	return ApplicationController{
		Config: config,
	}
}

func (c ApplicationController) Middleware(rw http.ResponseWriter, req *http.Request) {}
