package controllers

import (
	"encoding/json"
	"fmt"
	"golang-web-core/app/domain"
	"golang-web-core/app/models"
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/route"
	"golang-web-core/srv/srverr"
	"golang-web-core/util"
	"net/http"
	"reflect"
)

type AuthController struct {
	cfg.Config
}

func NewAuthController(cfg cfg.Config) AuthController {
	return AuthController{
		Config: cfg,
	}
}

// BeforeAction implements Controller.
func (a AuthController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
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
			Pattern:        "/auth/local/signup",
			Method:         http.MethodPost,
			Handler:        a.LocalSignUp,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/local/signup",
			Method:         http.MethodGet,
			Handler:        a.LocalSignUp,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/current_user",
			Method:         http.MethodGet,
			Handler:        a.CurrentUser,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/logout",
			Method:         http.MethodGet,
			Handler:        a.Logout,
			ControllerName: a.Name(),
		},
		{
			Pattern:        "/auth/logout",
			Method:         http.MethodDelete,
			Handler:        a.Logout,
			ControllerName: a.Name(),
		},
	}
}

func (a AuthController) LocalLogin(rw http.ResponseWriter, req *http.Request) {
	params, err := util.GetParams(req)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	email := ""
	password := ""

	if params["email"] != nil {
		email = params["email"].(string)
	}

	if params["password"] != nil {
		password = params["password"].(string)
	}

	if email == "" || password == "" {
		srverr.Handle400(rw, fmt.Errorf("email and password are required"))
		return
	}

	usersModel := models.NewUserModel(&a.Database.Adapter)
	user, err := usersModel.FindByEmail(email)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	passwordCorrect, err := user.CheckPassword(password)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	if !passwordCorrect {
		srverr.Handle401(rw, fmt.Errorf("password is incorrect"))
		return
	}

	sessionModel := models.NewSessionModel(&a.Database.Adapter)
	session, err := sessionModel.FindOrCreate(user.Id)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	if params["rememberMe"] == true {
		session.Expires = false
	}

	session.ResetExpiration()

	err = sessionModel.Update(session.Id.Hex(), session)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	if req.Method == http.MethodGet {
		http.SetCookie(rw, &http.Cookie{
			Name:  "session",
			Value: session.Token,
			Path:  "/",
		})
		http.Redirect(rw, req, "/auth/current_user", http.StatusTemporaryRedirect)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(session)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}
}

func (a AuthController) LocalSignUp(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(a.Database.Connection.Database)

	params, err := util.GetParams(req)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	email := ""
	password := ""
	firstName := ""
	lastName := ""

	if params["email"] != nil {
		email = params["email"].(string)
	}
	if params["password"] != nil {
		password = params["password"].(string)
	}
	if params["firstName"] != nil {
		firstName = params["firstName"].(string)
	}
	if params["lastName"] != nil {
		lastName = params["lastName"].(string)
	}

	if firstName == "" || lastName == "" {
		srverr.Handle400(rw, fmt.Errorf("first name and last name are required"))
		return
	}

	user, err := domain.NewUser(email, firstName, lastName, password)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	usersModel := models.NewUserModel(&a.Database.Adapter)
	_, err = usersModel.Create(user)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	a.LocalLogin(rw, req)
}

func (a AuthController) CurrentUser(rw http.ResponseWriter, req *http.Request) {
	user := util.GetContextUser(req)

	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(user)
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}
}

func (a AuthController) Logout(rw http.ResponseWriter, req *http.Request) {
	user := util.GetContextUser(req)
	if user == nil {
		if req.Method == http.MethodGet {
			http.Redirect(rw, req, "/auth/current_user", http.StatusTemporaryRedirect)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte(""))
		return
	}

	sessionModel := models.NewSessionModel(&a.Database.Adapter)
	err := sessionModel.DeleteWhere(map[string]any{"userId": user.Id})
	if err != nil {
		srverr.HandleSrvError(rw, err)
		return
	}

	if req.Method == http.MethodGet {
		http.Redirect(rw, req, "/auth/current_user", http.StatusTemporaryRedirect)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	rw.Write([]byte(""))
}

var _ Controller = AuthController{}
