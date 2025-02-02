package srv

import (
	"golang-web-core/app/controllers"
	"golang-web-core/srv/route"
	"log"
	"net/http"
)

func HandleRequest(appController controllers.ApplicationController, route route.Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		logRequest(req)

		controller := appController.Controllers[route.ControllerName]

		appController.BeforeAction(controller.BeforeAction(route.Handler))(rw, req)
	}
}

func HandleOptions(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.WriteHeader(http.StatusOK)
}

func logRequest(req *http.Request) {
	color := "255;255;255"

	switch req.Method {
	case http.MethodGet:
		color = "0;0;255"
	case http.MethodConnect:
		color = "0;0;255"
	case http.MethodOptions:
		color = "0;0;255"
	case http.MethodTrace:
		color = "0;0;255"
	case http.MethodPost:
		color = "100;255;100"
	case http.MethodPatch:
		color = "100;255;100"
	case http.MethodPut:
		color = "100;255;100"
	case http.MethodDelete:
		color = "255;0;0"
	}

	log.Printf("Started \033[38;2;%vm%v\033[0m %v for %v\n", color, req.Method, req.URL.Path, req.RemoteAddr)
}
