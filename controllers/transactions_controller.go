package controllers

import (
	"fmt"
	"inventory-app/domain"
	transactionrepo "inventory-app/repositories/transaction"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"net/http"
	"reflect"
)

type TransactionsController struct {
	transactionRepo domain.TransactionRepository
}

func NewTransactionsController(transactionRepo domain.TransactionRepository) TransactionsController {
	return TransactionsController{
		transactionRepo: transactionRepo,
	}
}

func NewTransactionsControllerFromConfig(config cfg.Config) (TransactionsController, error) {
	if !config.Mongo.IsEnabled() {
		return TransactionsController{}, fmt.Errorf("mongo is not enabled")
	}

	var transactionRepo domain.TransactionRepository

	switch config.TransactionRepository {
	case "MongoTransactionRepository":
		transactionRepo = transactionrepo.NewMongoTransactionRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return TransactionsController{}, fmt.Errorf("invalid transaction repository: %v", config.TransactionRepository)
	}

	return NewTransactionsController(transactionRepo), nil
}

// BeforeAction implements Controller.
func (t TransactionsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

// Name implements Controller.
func (t TransactionsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionsController) Routes() []route.Route {
	return []route.Route{}
}

var _ Controller = TransactionsController{}
