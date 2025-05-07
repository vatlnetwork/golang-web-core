package controllers

import (
	"inventory-app/domain"
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
