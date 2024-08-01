package srv

import (
	"context"
	"fmt"
	"golang-web-core/app/controllers"
	"golang-web-core/app/routes"
	"golang-web-core/srv/cfg"
	"golang-web-core/util"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Config cfg.Config
	Router routes.Router
	Mux    http.ServeMux
	Routes map[string]routes.Route
}

func NewServer(c cfg.Config) (*Server, error) {
	server := Server{
		Config: c,
		Router: routes.NewRouter(c),
		Mux:    *http.NewServeMux(),
		Routes: map[string]routes.Route{},
	}

	err := server.registerRoutes()
	if err != nil {
		return &server, err
	}

	fmt.Println("Server Config:")
	fmt.Printf("   Port: %v\n", c.Port)
	fmt.Printf("   Using SSL: %v\n", c.IsSSL())
	if c.IsSSL() {
		fmt.Printf("      Cert Path: %v\n", c.SSL.CertPath)
		fmt.Printf("      Key Path: %v\n", c.SSL.KeyPath)
	}
	fmt.Printf("   # of Routes: %v\n", len(server.Routes))
	fmt.Println("")

	return &server, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Config.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server := http.Server{
		Addr:    addr,
		Handler: &s.Mux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				log.Println("Gracefully shutting down...")
				err := server.Shutdown(ctx)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	util.LogColor("green", "Server listening on %v\n", addr)

	if s.Config.IsSSL() {
		return server.ServeTLS(l, s.Config.SSL.CertPath, s.Config.SSL.KeyPath)
	}
	return server.Serve(l)
}

func (s *Server) registerRoutes() error {
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
	return nil
}
