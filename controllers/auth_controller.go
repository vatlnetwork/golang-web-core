package controllers

import (
	"fmt"
	"inventory-app/domain"
	sessionrepo "inventory-app/repositories/session"
	userrepo "inventory-app/repositories/user"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"net/http"
	"reflect"
)

type AuthController struct {
	sessionManager domain.SessionManager
}

func NewAuthController(sessionManager domain.SessionManager) AuthController {
	return AuthController{
		sessionManager: sessionManager,
	}
}

func NewAuthControllerFromConfig(config cfg.Config) (AuthController, error) {
	if !config.Mongo.IsEnabled() {
		return AuthController{}, fmt.Errorf("mongo is not enabled")
	}

	var userRepo domain.UserRepository
	var sessionRepo domain.SessionRepository

	switch config.UserRepository {
	case "MongoUserRepository":
		userRepo = userrepo.NewMongoUserRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return AuthController{}, fmt.Errorf("invalid user repository: %v", config.UserRepository)
	}

	switch config.SessionRepository {
	case "MongoSessionRepository":
		sessionRepo = sessionrepo.NewMongoSessionRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return AuthController{}, fmt.Errorf("invalid session repository: %v", config.SessionRepository)
	}

	sessionManager := domain.NewSessionManager(sessionRepo, userRepo)

	return NewAuthController(sessionManager), nil
}

// BeforeAction implements Controller.
func (a AuthController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

// Name implements Controller.
func (a AuthController) Name() string {
	return reflect.TypeOf(a).Name()
}

func (a AuthController) Routes() []route.Route {
	return []route.Route{}
}

func (a AuthController) LocalLogin(rw http.ResponseWriter, req *http.Request) {}

func (a AuthController) Logout(rw http.ResponseWriter, req *http.Request) {}

func (a AuthController) CurrentUser(rw http.ResponseWriter, req *http.Request) {}

var _ Controller = AuthController{}
