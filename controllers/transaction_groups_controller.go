package controllers

import (
	"inventory-app/domain"
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
