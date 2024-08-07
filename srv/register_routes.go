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
		_, ok := s.Routes[route.Pattern]
		if ok {
			return fmt.Errorf("error: route pattern %v was registered twice. you may only register a single pattern once", route.Pattern)
		}
		s.Routes[route.Pattern] = route

		s.Mux.HandleFunc(fmt.Sprintf("%v %v", route.Method, route.Pattern), HandleRequest(appController, route))
		s.Mux.HandleFunc(fmt.Sprintf("%v %v", http.MethodOptions, route.Pattern), http.HandlerFunc(HandleOptions))
	}

	if s.Config.PublicFS {
		s.Mux.Handle("/public/", http.StripPrefix("/public/", FileServer{Handler: http.FileServer(http.Dir("public"))}))
	}

	return nil
}
