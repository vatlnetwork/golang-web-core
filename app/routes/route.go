package routes

import "net/http"

type Route struct {
	Pattern string
	Method  int
	Handler http.HandlerFunc
}
