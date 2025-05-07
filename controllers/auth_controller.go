package controllers

import (
	"encoding/json"
	"fmt"
	"inventory-app/domain"
	sessionrepo "inventory-app/repositories/session"
	userrepo "inventory-app/repositories/user"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"inventory-app/srv/srverr"
	"inventory-app/util"
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
	var userRepo domain.UserRepository
	var sessionRepo domain.SessionRepository

	switch config.UserRepository {
	case "MongoUserRepository":
		if !config.Mongo.IsEnabled() {
			return AuthController{}, fmt.Errorf("mongo is not enabled")
		}
		userRepo = userrepo.NewMongoUserRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return AuthController{}, fmt.Errorf("invalid user repository: %v", config.UserRepository)
	}

	switch config.SessionRepository {
	case "MongoSessionRepository":
		if !config.Mongo.IsEnabled() {
			return AuthController{}, fmt.Errorf("mongo is not enabled")
		}
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
	return []route.Route{
		{
			Pattern:        "/auth/local/login",
			Method:         http.MethodPost,
			Handler:        a.LocalLogin,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/local/login",
			Method:         http.MethodGet,
			Handler:        a.LocalLogin,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/logout",
			Method:         http.MethodDelete,
			Handler:        a.Logout,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/logout",
			Method:         http.MethodGet,
			Handler:        a.Logout,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/current_user",
			Method:         http.MethodGet,
			Handler:        a.CurrentUser,
			ControllerName: a.Name(),
		},
	}
}

type localLoginRequest struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Password     string `json:"password"`
	StaySignedIn string `json:"staySignedIn"`
}

type localLoginResponse struct {
	Session domain.Session `json:"session"`
	User    domain.User    `json:"user"`
}

func (a AuthController) LocalLogin(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	currentSession, currentUser, err := a.sessionManager.GetCurrentSession(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	if currentSession != nil {
		response := localLoginResponse{
			Session: *currentSession,
			User:    currentUser,
		}

		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		return
	}

	var request localLoginRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	session, user, err := a.sessionManager.HandleSignIn(
		req,
		request.Email,
		request.FirstName,
		request.LastName,
		request.Password,
		request.StaySignedIn == "yes",
	)
	if err != nil {
		if err.Error() == domain.ErrorInvalidEmail || err.Error() == domain.ErrorInvalidPassword {
			srverr.Handle400(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if req.Method == http.MethodGet {
		cookie := &http.Cookie{
			Name:  "session",
			Value: session.Id,
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
		http.Redirect(rw, req, "/auth/current_user", http.StatusSeeOther)
	} else {
		response := localLoginResponse{
			Session: session,
			User:    user,
		}

		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}
	}
}

func (a AuthController) Logout(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := a.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	if currentUser != nil {
		err = a.sessionManager.HandleSignOut(currentUser.Id)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}
	}

	if req.Method == http.MethodGet {
		http.Redirect(rw, req, "/auth/current_user", http.StatusSeeOther)
	} else {
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (a AuthController) CurrentUser(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := a.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(currentUser)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

var _ Controller = AuthController{}
