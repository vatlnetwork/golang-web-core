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
	// this is where you define your routes. you can do this however you like so long as you populate
	// all of the fields in each route. each field is necessary. if you have a lot of routes, you can split
	// your routes up into multiple files, so long as they are all returned here
	testController := appController.GetController("TestController").(controllers.TestController)

	return []Route{
		{
			Pattern:        "/test_route",
			Method:         http.MethodGet,
			Handler:        testController.TestMethod,
			ControllerName: testController.Name(),
		},
		{
			Pattern:        "/test_member_route/{member_var}/test",
			Method:         http.MethodGet,
			Handler:        testController.TestMemberMethod,
			ControllerName: testController.Name(),
		},
	}
}
