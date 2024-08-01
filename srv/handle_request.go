package srv

import (
	"golang-web-core/app/controllers"
	"golang-web-core/app/routes"
	"net/http"
)

func HandleRequest(appController controllers.ApplicationController, route routes.Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		appController.BeforeAction(route.Controller.BeforeAction(route.Handler))
	}
}

func HandleOptions(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.WriteHeader(http.StatusOK)
}
