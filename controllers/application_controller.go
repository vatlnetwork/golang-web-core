package controllers

import (
	"fmt"
	"inventory-app/domain"
	moneylocationrepo "inventory-app/repositories/money_location"
	sessionrepo "inventory-app/repositories/session"
	transactionrepo "inventory-app/repositories/transaction"
	transactiongrouprepo "inventory-app/repositories/transaction_group"
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
	Controllers          map[string]Controller
	sessionManager       domain.SessionManager
	userRepo             domain.UserRepository
	sessionRepo          domain.SessionRepository
	transactionRepo      domain.TransactionRepository
	transactionGroupRepo domain.TransactionGroupRepository
	moneyLocationRepo    domain.MoneyLocationRepository
}

// this verifies that ApplicationController fully implements Controller
var ApplicationControllerVerifier Controller = ApplicationController{}

func NewApplicationController(config cfg.Config) (ApplicationController, error) {
	cont := ApplicationController{
		Config:      config,
		Controllers: map[string]Controller{},
	}

	err := cont.setupRepositories()
	if err != nil {
		return ApplicationController{}, err
	}

	cont.sessionManager = domain.NewSessionManager(cont.sessionRepo, cont.userRepo)

	err = cont.setupControllers()
	if err != nil {
		return ApplicationController{}, err
	}

	return cont, nil
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

func (c *ApplicationController) setupRepositories() error {
	var transactionRepo domain.TransactionRepository
	switch c.Config.TransactionRepository {
	case "MongoTransactionRepository":
		if !c.Config.Mongo.IsEnabled() {
			return fmt.Errorf("mongo is not enabled")
		}
		transactionRepo = transactionrepo.NewMongoTransactionRepository(c.Config.Mongo, c.Config.Env == cfg.Development)
	default:
		return fmt.Errorf("invalid transaction repository: %v", c.Config.TransactionRepository)
	}
	c.transactionRepo = transactionRepo

	var transactionGroupRepo domain.TransactionGroupRepository
	switch c.Config.TransactionGroupRepository {
	case "MongoTransactionGroupRepository":
		if !c.Config.Mongo.IsEnabled() {
			return fmt.Errorf("mongo is not enabled")
		}
		transactionGroupRepo = transactiongrouprepo.NewMongoTransactionGroupRepository(c.Config.Mongo, c.Config.Env == cfg.Development)
	default:
		return fmt.Errorf("invalid transaction group repository: %v", c.Config.TransactionGroupRepository)
	}
	c.transactionGroupRepo = transactionGroupRepo

	var userRepo domain.UserRepository
	switch c.Config.UserRepository {
	case "MongoUserRepository":
		if !c.Config.Mongo.IsEnabled() {
			return fmt.Errorf("mongo is not enabled")
		}
		userRepo = userrepo.NewMongoUserRepository(c.Config.Mongo, c.Config.Env == cfg.Development)
	default:
		return fmt.Errorf("invalid user repository: %v", c.Config.UserRepository)
	}
	c.userRepo = userRepo

	var sessionRepo domain.SessionRepository
	switch c.Config.SessionRepository {
	case "MongoSessionRepository":
		if !c.Config.Mongo.IsEnabled() {
			return fmt.Errorf("mongo is not enabled")
		}
		sessionRepo = sessionrepo.NewMongoSessionRepository(c.Config.Mongo, c.Config.Env == cfg.Development)
	default:
		return fmt.Errorf("invalid session repository: %v", c.Config.SessionRepository)
	}
	c.sessionRepo = sessionRepo

	var moneyLocationRepo domain.MoneyLocationRepository
	switch c.Config.MoneyLocationRepository {
	case "MongoMoneyLocationRepository":
		if !c.Config.Mongo.IsEnabled() {
			return fmt.Errorf("mongo is not enabled")
		}
		moneyLocationRepo = moneylocationrepo.NewMongoMoneyLocationRepository(c.Config.Mongo, c.Config.Env == cfg.Development)
	default:
		return fmt.Errorf("invalid money location repository: %v", c.Config.MoneyLocationRepository)
	}
	c.moneyLocationRepo = moneyLocationRepo

	return nil
}

func (c *ApplicationController) setupControllers() error {
	controllers := []Controller{
		c,
		// this is where you initialize your controllers. if you do not initialize your controllers here, they will not be usable
		NewTransactionsController(c.transactionRepo, c.sessionManager, c.transactionGroupRepo, c.moneyLocationRepo),
		NewTransactionGroupsController(c.transactionGroupRepo, c.sessionManager),
		NewAuthController(c.sessionManager),
		NewMoneyLocationsController(c.moneyLocationRepo, c.sessionManager),
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
