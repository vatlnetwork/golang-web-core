package controllers

import (
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/route"
	"net/http"
	"reflect"
)

type TransactionsController struct {
	cfg.Config
}

// this verifies that TransactionsController fully implements Controller
var TransactionsControllerVerifier Controller = TransactionsController{}

func NewTransactionsController(c cfg.Config) TransactionsController {
	return TransactionsController{
		Config: c,
	}
}

func (c TransactionsController) Name() string {
	return reflect.TypeOf(c).Name()
}

func (c TransactionsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

func (c TransactionsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/transactions/{year}",
			Method:         http.MethodGet,
			Handler:        c.Load,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/transactions",
			Method:         http.MethodPost,
			Handler:        c.Create,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/transactions/{id}/update",
			Method:         http.MethodPatch,
			Handler:        c.Update,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/transactions/{id}/delete",
			Method:         http.MethodDelete,
			Handler:        c.Delete,
			ControllerName: c.Name(),
		},
	}
}

func (c TransactionsController) Load(rw http.ResponseWriter, req *http.Request) {}

func (c TransactionsController) Create(rw http.ResponseWriter, req *http.Request) {}

func (c TransactionsController) Update(rw http.ResponseWriter, req *http.Request) {}

func (c TransactionsController) Delete(rw http.ResponseWriter, req *http.Request) {}
