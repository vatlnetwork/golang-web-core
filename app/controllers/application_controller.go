package controllers

import (
	"context"
	"fmt"
	"golang-web-core/app/models"
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/srverr"
	"golang-web-core/util"
	"net/http"
	"reflect"
)

// you shouldn't be touching this file except for the BeforeAction and setupControllers

type ApplicationController struct {
	cfg.Config
	Controllers map[string]Controller
}

// this verifies that ApplicationController fully implements Controller
var ApplicationControllerVerifier Controller = ApplicationController{}

func NewApplicationController(config cfg.Config) (ApplicationController, error) {
	cont := ApplicationController{
		Config:      config,
		Controllers: map[string]Controller{},
	}

	err := cont.setupControllers()

	return cont, err
}

func (c ApplicationController) Name() string {
	return reflect.TypeOf(c).Name()
}

func (c ApplicationController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// any checks you want to do on every single request that goes into the server can go here

		// this line initiates the next step in the request process.
		// if you wanted to throw an error here or something, it might look something like this:
		// if (someCondition) {
		// 	http.Error(rw, "This is a test internal server error", http.StatusInternalServerError)
		// 	return
		// }

		token := req.Header.Get("Authorization")
		if token == "" {
			sessionCookie, err := req.Cookie("session")
			if err != nil {
				srverr.HandleSrvError(rw, err)
				return
			}

			token = sessionCookie.Value

			if token == "" {
				handler(rw, req)
				return
			}
		}

		sessionModel := models.NewSessionModel(&c.Database.Adapter)
		session, err := sessionModel.FindByToken(token)
		if err != nil {
			if err.(srverr.ServerError).Code != http.StatusNotFound {
				srverr.HandleSrvError(rw, err)
				return
			}
		}

		if session == nil {
			handler(rw, req)
			return
		}

		if !session.IsExpired() {
			userModel := models.NewUserModel(&c.Database.Adapter)
			user, err := userModel.Find(session.UserId.Hex())
			if err != nil {
				srverr.HandleSrvError(rw, err)
				return
			}

			reqWithUser := req.WithContext(context.WithValue(req.Context(), "current_user", user))
			handler(rw, reqWithUser)
			return
		}

		handler(rw, req)
	}
}

func (c ApplicationController) Favicon(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "app/favicon.ico")
}

func (c ApplicationController) setupControllers() error {
	controllers := []Controller{
		c,
		// this is where you initialize your controllers. if you do not initialize your controllers here, they will not be usable
		NewTestController(c.Config),
		NewTransactionsController(c.Config),
		NewTransactionGroupsController(c.Config),
		NewAuthController(c.Config),
	}

	// everything below here should be left untouched

	for _, cont := range controllers {
		_, ok := c.Controllers[cont.Name()]
		if ok {
			return fmt.Errorf("error: a controller with the name %v was registered twice", cont.Name())
		}
		c.Controllers[cont.Name()] = cont
	}

	return nil
}

func (c ApplicationController) GetController(name string) Controller {
	controller, ok := c.Controllers[name]
	if !ok {
		util.LogFatalf("attempted to access a controller that does not exist! %v; please add it to ApplicationController.setupControllers", name)
	}

	return controller
}
