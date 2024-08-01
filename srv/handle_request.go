package srv

import (
	"fmt"
	"golang-web-core/app/controllers"
	"golang-web-core/app/routes"
	"net/http"
)

func HandleRequest(appController controllers.ApplicationController, route routes.Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("this ran 1")

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
