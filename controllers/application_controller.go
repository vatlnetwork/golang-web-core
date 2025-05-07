package controllers

import (
	"fmt"
	"inventory-app/domain"
	sessionrepo "inventory-app/repositories/session"
	userrepo "inventory-app/repositories/user"
	"inventory-app/srv/cfg"
	"inventory-app/srv/srverr"
	"inventory-app/util"
	"net/http"
	"reflect"
)

// you shouldn't be touching this file except for the BeforeAction and setupControllers

type ApplicationController struct {
	cfg.Config
	Controllers    map[string]Controller
	sessionManager domain.SessionManager
}

// this verifies that ApplicationController fully implements Controller
var ApplicationControllerVerifier Controller = ApplicationController{}

func NewApplicationController(config cfg.Config) (ApplicationController, error) {
	var userRepo domain.UserRepository
	var sessionRepo domain.SessionRepository

	switch config.UserRepository {
	case "MongoUserRepository":
		if !config.Mongo.IsEnabled() {
			return ApplicationController{}, fmt.Errorf("mongo is not enabled")
		}
		userRepo = userrepo.NewMongoUserRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return ApplicationController{}, fmt.Errorf("invalid user repository: %v", config.UserRepository)
	}

	switch config.SessionRepository {
	case "MongoSessionRepository":
		if !config.Mongo.IsEnabled() {
			return ApplicationController{}, fmt.Errorf("mongo is not enabled")
		}
		sessionRepo = sessionrepo.NewMongoSessionRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return ApplicationController{}, fmt.Errorf("invalid session repository: %v", config.SessionRepository)
	}

	sessionManager := domain.NewSessionManager(sessionRepo, userRepo)

	cont := ApplicationController{
		Config:         config,
		Controllers:    map[string]Controller{},
		sessionManager: sessionManager,
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

		session, user, err := c.sessionManager.GetCurrentSession(req)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		if session != nil {
			req = c.sessionManager.SetContextUser(req, user)
		}

		handler(rw, req)
	}
}

func (c ApplicationController) Favicon(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "favicon.ico")
}

func (c ApplicationController) setupControllers() error {
	transactionsController, err := NewTransactionsControllerFromConfig(c.Config)
	if err != nil {
		return err
	}

	transactionGroupsController, err := NewTransactionGroupsControllerFromConfig(c.Config)
	if err != nil {
		return err
	}

	authController, err := NewAuthControllerFromConfig(c.Config)
	if err != nil {
		return err
	}

	moneyLocationsController, err := NewMoneyLocationsControllerFromConfig(c.Config)
	if err != nil {
		return err
	}

	controllers := []Controller{
		c,
		transactionsController,
		transactionGroupsController,
		authController,
		moneyLocationsController,
		// this is where you initialize your controllers. if you do not initialize your controllers here, they will not be usable
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
