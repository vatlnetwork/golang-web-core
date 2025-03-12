package srv

import (
	"fmt"
	"golang-web-core/app/controllers"
	"net/http"
)

func (s *Server) RegisterRoutes() error {
	appController, err := controllers.NewApplicationController(s.Config)
	if err != nil {
		return err
	}
	routes := s.Router.Routes(appController)

	for _, route := range routes {
		existingRoute, ok := s.Routes[route.Pattern]
		if ok {
			if existingRoute.Method == route.Method {
				return fmt.Errorf("error: route pattern %v %v was registered twice", route.Method, route.Pattern)
			}
		}
		s.Routes[route.Pattern] = route

		s.Mux.HandleFunc(fmt.Sprintf("%v %v", route.Method, route.Pattern), HandleRequest(appController, route))
		if !ok {
			s.Mux.HandleFunc(fmt.Sprintf("%v %v", http.MethodOptions, route.Pattern), http.HandlerFunc(HandleOptions))
		}
	}

	if s.Config.PublicFS {
		s.Mux.Handle("/public/", http.StripPrefix("/public/", FileServer{Prefix: "/public/", Handler: http.FileServer(http.Dir("public"))}))
	}

	return nil
}
