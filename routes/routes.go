package routes

import (
	"fmt"
	"net/http"
	"golang-web-core/controllers"
	"golang-web-core/services/httpserver"
	"slices"
)

func Routes(controllers []controllers.Controller, applicationController controllers.ApplicationController) ([]httpserver.Route, error) {
	routes := []httpserver.Route{}

	for _, route := range applicationController.Routes() {
		routes = append(routes, httpserver.Route{
			Pattern: route.Pattern,
			Method:  route.Method,
			Handler: applicationController.BeforeAction(route.Handler),
		})
	}

	for _, controller := range controllers {
		for _, route := range controller.Routes() {
			routes = append(routes, httpserver.Route{
				Pattern: route.Pattern,
				Method:  route.Method,
				Handler: applicationController.BeforeAction(controller.BeforeAction(route.Handler)),
			})
		}
	}

	patterns := []string{}
	registeredRoutes := map[string]httpserver.Route{}

	for _, route := range routes {
		_, ok := registeredRoutes[route.Method+" "+route.Pattern]
		if ok {
			return nil, fmt.Errorf("route pattern %v %v was registered twice", route.Method, route.Pattern)
		}
		registeredRoutes[route.Method+" "+route.Pattern] = route

		if !slices.Contains(patterns, route.Pattern) {
			patterns = append(patterns, route.Pattern)
			optionsRoute := httpserver.Route{
				Pattern: route.Pattern,
				Method:  http.MethodOptions,
				// Options routes don't require a handler because they are handled by the server
			}
			routes = append(routes, optionsRoute)
		}
	}

	return routes, nil
}
