package srv

import (
	"golang-web-core/app/controllers"
	"golang-web-core/app/routes"
	"log"
	"net/http"
)

func HandleRequest(appController controllers.ApplicationController, route routes.Route) http.HandlerFunc {
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
	log.Printf("Started %v %v\n", req.Method, req.URL.Path)
}
