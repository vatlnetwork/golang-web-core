package controllers

import (
	"golang-web-core/srv"
	"net/http"
)

type ApplicationController struct {
	srv.Config
}

func (c ApplicationController) Middleware(rw http.ResponseWriter, req *http.Request) {}
