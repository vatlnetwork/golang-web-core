package routes

import (
	"golang-web-core/app/controllers"
	"golang-web-core/srv/cfg"
	"net/http"
)

type Router struct {
	config cfg.Config
}

func NewRouter(c cfg.Config) Router {
	return Router{
		config: c,
	}
}

func (r Router) Routes(appController controllers.ApplicationController) []Route {
	testController := appController.GetController("TestController").(controllers.TestController)

	return []Route{
		{
			Pattern:        "/test_route",
			Method:         http.MethodGet,
			Handler:        testController.TestMethod,
			ControllerName: "TestController",
		},
	}
}
