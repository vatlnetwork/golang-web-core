package controllers

import (
	"golang-web-core/srv/cfg"
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
