package controllers

import (
	"fmt"
	"inventory-app/domain"
	transactiongrouprepo "inventory-app/repositories/transaction_group"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"net/http"
	"reflect"
)

type TransactionGroupsController struct {
	transactionGroupRepo domain.TransactionGroupRepository
}

func NewTransactionGroupsController(transactionGroupRepo domain.TransactionGroupRepository) TransactionGroupsController {
	return TransactionGroupsController{
		transactionGroupRepo: transactionGroupRepo,
	}
}

func NewTransactionGroupsControllerFromConfig(config cfg.Config) (TransactionGroupsController, error) {
	if !config.Mongo.IsEnabled() {
		return TransactionGroupsController{}, fmt.Errorf("mongo is not enabled")
	}

	var transactionGroupRepo domain.TransactionGroupRepository

	switch config.TransactionGroupRepository {
	case "MongoTransactionGroupRepository":
		transactionGroupRepo = transactiongrouprepo.NewMongoTransactionGroupRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return TransactionGroupsController{}, fmt.Errorf("invalid transaction group repository: %v", config.TransactionGroupRepository)
	}

	return NewTransactionGroupsController(transactionGroupRepo), nil
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
