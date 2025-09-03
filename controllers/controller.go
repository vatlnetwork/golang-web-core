package controllers

import (
	"net/http"
	"golang-web-core/services/httpserver"
)

type Controller interface {
	BeforeAction(handler http.HandlerFunc) http.HandlerFunc
	Routes() []httpserver.Route
}
