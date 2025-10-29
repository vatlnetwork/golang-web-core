package controllers

import (
	"encoding/json"
	"errors"
	"golang-web-core/domain"
	"golang-web-core/services/httpserver"
	"net/http"
	"time"
)

type AuthController struct {
	sessionManager domain.SessionManager
	errorHandler   *httpserver.HttpErrorHandler
}

func NewAuthController(sessionManager domain.SessionManager, errorHandler *httpserver.HttpErrorHandler) (AuthController, error) {
	if errorHandler == nil {
		return AuthController{}, errors.New("error handler is required")
	}

	return AuthController{
		sessionManager: sessionManager,
		errorHandler:   errorHandler,
	}, nil
}

// BeforeAction implements httpserver.Controller.
func (a AuthController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

// Routes implements httpserver.Controller.
func (a AuthController) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Pattern: "/auth/local/login",
			Method:  http.MethodPost,
			Handler: a.LocalLogin,
		},
		{
			Pattern: "/auth/local/logout",
			Method:  http.MethodDelete,
			Handler: a.LocalLogout,
		},
		{
			Pattern: "/auth/local/logout",
			Method:  http.MethodGet,
			Handler: a.LocalLogout,
		},
		{
			Pattern: "/auth/current_user",
			Method:  http.MethodGet,
			Handler: a.CurrentUser,
		},
	}
}

type localLoginRequest struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Password     string `json:"password"`
	StaySignedIn bool   `json:"staySignedIn"`
}

type localLoginResponse struct {
	Session domain.Session `json:"session"`
	User    domain.User    `json:"user"`
}

func (a AuthController) LocalLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currentSession, currentUser, err := a.sessionManager.GetCurrentSession(r)
	if err != nil {
		a.errorHandler.HandleError(http.StatusUnauthorized, w, err)
		return
	}

	if currentSession != nil {
		response := localLoginResponse{
			Session: *currentSession,
			User:    currentUser,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
			return
		}

		return
	}

	request := localLoginRequest{}
	err = httpserver.DecodeContextParams(r, &request)
	if err != nil {
		a.errorHandler.HandleError(http.StatusBadRequest, w, err)
		return
	}

	session, user, err := a.sessionManager.HandleSignIn(
		r,
		request.Email,
		request.FirstName,
		request.LastName,
		request.Password,
		request.StaySignedIn,
	)
	if err != nil {
		if err.Error() == domain.ErrorInvalidEmail || err.Error() == domain.ErrorInvalidPassword {
			a.errorHandler.HandleError(http.StatusBadRequest, w, err)
		} else {
			a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
		}
		return
	}

	cookie := &http.Cookie{
		Name:  domain.SessionCookieName,
		Value: session.Id,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	response := localLoginResponse{
		Session: session,
		User:    user,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
		return
	}
}

func (a AuthController) LocalLogout(w http.ResponseWriter, r *http.Request) {
	currentSession, _, err := a.sessionManager.GetCurrentSession(r)
	if err != nil {
		a.errorHandler.HandleError(http.StatusUnauthorized, w, err)
		return
	}

	if currentSession != nil {
		err = a.sessionManager.HandleSignOut(currentSession.Id)
		if err != nil {
			a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
			return
		}

		cookie := &http.Cookie{
			Name:    domain.SessionCookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-time.Hour * 24),
		}
		http.SetCookie(w, cookie)
	}

	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/auth/current_user", http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a AuthController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUser, err := a.sessionManager.GetContextUser(r)
	if err != nil {
		a.errorHandler.HandleError(http.StatusUnauthorized, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(currentUser)
	if err != nil {
		a.errorHandler.HandleError(http.StatusInternalServerError, w, err)
		return
	}
}

var _ httpserver.Controller = AuthController{}
