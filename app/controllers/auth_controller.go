package controllers

import (
	"fmt"
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
	return []route.Route{}
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

	// usersModel := models.NewUserModel(&a.Database.Adapter)
}

var _ Controller = AuthController{}
