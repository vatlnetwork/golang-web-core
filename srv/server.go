package srv

import (
	"fmt"
	"golang-web-core/app/routes"
	"golang-web-core/srv/cfg"
	"golang-web-core/util"
	"net"
	"net/http"
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

	err := server.RegisterRoutes()
	if err != nil {
		return &server, err
	}

	PrintServerConfig(&server)

	return &server, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Config.Port)
	server := http.Server{
		Addr:    addr,
		Handler: &s.Mux,
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if s.Config.Database.Adapter != nil {
		err = s.TestDatabase()
		if err != nil {
			return err
		}
	}

	s.HandleShutdown(&server)

	util.LogColor("lightgreen", "Server listening on %v\n", addr)

	if s.Config.IsSSL() {
		return server.ServeTLS(l, s.Config.SSL.CertPath, s.Config.SSL.KeyPath)
	}
	return server.Serve(l)
}
