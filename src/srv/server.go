package srv

import (
	"fmt"
	"golang-web-core/src/interface/database/mongo"
	"golang-web-core/src/srv/cfg"
	"golang-web-core/src/srv/middlewares"
	"golang-web-core/src/srv/routes"
	"golang-web-core/src/util"
	"log"
	"net/http"
	"runtime"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

type Server struct {
	Config cfg.ServerConfig
	*mux.Router
}

func NewServer(config cfg.ServerConfig) (*Server, error) {
	s := &Server{
		Router: mux.NewRouter(),
		Config: config,
	}
	err := s.Configure()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func options(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.WriteHeader(http.StatusOK)
}

func (s *Server) Configure() error {
	// gather the routes from the controllers
	srvRoutes, err := acquireRoutes(s)
	if err != nil {
		return err
	}

	// map the routes to their patterns
	routesMap := map[string]routes.Route{}
	for _, route := range srvRoutes {
		routesMap[route.Path] = route
	}

	// register the routes with mux router
	for _, route := range srvRoutes {
		s.HandleFunc(route.Path, route.Handler).Methods(route.Method)
		s.HandleFunc(route.Path, options).Methods(http.MethodOptions)
	}

	// setup frontend application routes
	appDir, err := util.AppDir()
	if err != nil {
		return err
	}
	frontendDir := "/resources/frontend"
	if runtime.GOOS == "windows" {
		frontendDir = "\\resources\\frontend"
	}
	s.PathPrefix("/").Handler(http.FileServer(http.Dir(fmt.Sprintf("%v%v", appDir, frontendDir))))

	// register middleware functions
	s.Use(middlewares.EnableCors)
	s.Use(middlewares.RequestLogger)
	s.Use(middlewares.AuthMiddleware(routesMap, s.Router, s.Repos().Sessions))

	// output server configuration / stats
	fmt.Printf("Server Port:          %v\n", s.Config.Port)
	fmt.Printf("Using SSL:            %v\n", s.Config.UsesSSL())
	fmt.Printf("MongoDB Server Host:  %v\n", s.Config.Mongo.Host)
	fmt.Printf("MongoDB Server Port:  %v\n", s.Config.Mongo.Port)
	fmt.Printf("MongoDB Database:     %v\n", s.Config.Mongo.DbName)
	fmt.Printf("Using MongoDB Auth:   %v\n", s.Config.Mongo.UsesAuth())
	fmt.Printf("Media Directory:      %v\n", s.Config.Media.Directory)
	fmt.Printf("# of Routes:          %v\n\n", len(srvRoutes))

	log.Println("Testing DB Connection...")
	if s.Config.Env != cfg.MockDev {
		client, context, cancelFunc, err := mongo.Connect(s.Config.Mongo.ConnectionString())
		if err != nil {
			return err
		}
		defer mongo.Close(client, context, cancelFunc)
		err = mongo.Ping(client, context)
		if err != nil {
			return err
		}
		log.Println("DB Connection Good")
	} else {
		log.Println("No DB required in mock development mode")
	}

	return nil
}

func (s Server) Run() error {
	c := color.New(color.FgGreen)
	c.Printf("Server initialized on port %v\n", s.Config.Port)
	var err error
	if s.Config.UsesSSL() {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%v", s.Config.Port), s.Config.SSL.CertFilePath, s.Config.SSL.CertKeyFilePath, s)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), s)
	}
	if err != nil {
		return err
	}
	return nil
}

func acquireRoutes(s *Server) ([]routes.Route, error) {
	// setup routes
	srvRoutes := []routes.Route{}
	srvRoutes = append(srvRoutes, s.Controllers().Users.Routes()...)
	srvRoutes = append(srvRoutes, s.Controllers().Sessions.Routes()...)
	srvRoutes = append(srvRoutes, s.Controllers().Media.Routes()...)

	// check routes
	paths := []string{}
	for _, route := range srvRoutes {
		if util.ArrayContainsString(paths, route.Path) {
			return []routes.Route{}, fmt.Errorf("error: there are 2 or more routes with the same path: %v", route.Path)
		}
		paths = append(paths, route.Path)
	}

	return srvRoutes, nil
}
