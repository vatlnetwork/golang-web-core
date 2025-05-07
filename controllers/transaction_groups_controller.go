package controllers

import (
	"fmt"
	"inventory-app/domain"
	sessionrepo "inventory-app/repositories/session"
	transactiongrouprepo "inventory-app/repositories/transaction_group"
	userrepo "inventory-app/repositories/user"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"net/http"
	"reflect"
)

type TransactionGroupsController struct {
	transactionGroupRepo domain.TransactionGroupRepository
	sessionManager       domain.SessionManager
}

func NewTransactionGroupsController(
	transactionGroupRepo domain.TransactionGroupRepository,
	sessionManager domain.SessionManager,
) TransactionGroupsController {
	return TransactionGroupsController{
		transactionGroupRepo: transactionGroupRepo,
		sessionManager:       sessionManager,
	}
}

func NewTransactionGroupsControllerFromConfig(config cfg.Config) (TransactionGroupsController, error) {
	var transactionGroupRepo domain.TransactionGroupRepository
	var usersRepo domain.UserRepository
	var sessionsRepo domain.SessionRepository

	switch config.TransactionGroupRepository {
	case "MongoTransactionGroupRepository":
		if !config.Mongo.IsEnabled() {
			return TransactionGroupsController{}, fmt.Errorf("mongo is not enabled")
		}
		transactionGroupRepo = transactiongrouprepo.NewMongoTransactionGroupRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return TransactionGroupsController{}, fmt.Errorf("invalid transaction group repository: %v", config.TransactionGroupRepository)
	}

	switch config.UserRepository {
	case "MongoUserRepository":
		if !config.Mongo.IsEnabled() {
			return TransactionGroupsController{}, fmt.Errorf("mongo is not enabled")
		}
		usersRepo = userrepo.NewMongoUserRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return TransactionGroupsController{}, fmt.Errorf("invalid user repository: %v", config.UserRepository)
	}

	switch config.SessionRepository {
	case "MongoSessionRepository":
		if !config.Mongo.IsEnabled() {
			return TransactionGroupsController{}, fmt.Errorf("mongo is not enabled")
		}
		sessionsRepo = sessionrepo.NewMongoSessionRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return TransactionGroupsController{}, fmt.Errorf("invalid session repository: %v", config.SessionRepository)
	}

	sessionManager := domain.NewSessionManager(sessionsRepo, usersRepo)

	return NewTransactionGroupsController(transactionGroupRepo, sessionManager), nil
}

// BeforeAction implements Controller.
func (t TransactionGroupsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

// Name implements Controller.
func (t TransactionGroupsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionGroupsController) Routes() []route.Route {
	return []route.Route{}
}

var _ Controller = TransactionGroupsController{}
