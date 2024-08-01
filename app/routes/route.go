package routes

import (
	"golang-web-core/app/controllers"
	"net/http"
)

type Route struct {
	Pattern    string
	Method     int
	Handler    http.HandlerFunc
	Controller controllers.Controller
}
