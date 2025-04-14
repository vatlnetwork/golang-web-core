package controllers

import (
	"encoding/json"
	"fmt"
	"golang-web-core/app/domain"
	"golang-web-core/app/models"
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/route"
	"golang-web-core/srv/srverr"
	"golang-web-core/util"
	"net/http"
	"reflect"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionsController struct {
	cfg.Config
}

func NewTransactionsController(cfg cfg.Config) TransactionsController {
	return TransactionsController{
		Config: cfg,
	}
}

// BeforeAction implements Controller.
func (t TransactionsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := util.ContextUserOrError(r)
		if err != nil {
			srverr.Handle401(w, err)
			return
		}

		handler(w, r)
	}
}

// Name implements Controller.
func (t TransactionsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/transactions",
			Method:         http.MethodPost,
			Handler:        t.CreateTransaction,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transactions/year/{year}",
			Method:         http.MethodGet,
			Handler:        t.GetTransactionsByYear,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transactions/{id}",
			Method:         http.MethodPut,
			Handler:        t.UpdateTransaction,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transactions/{id}",
			Method:         http.MethodDelete,
			Handler:        t.DeleteTransaction,
			ControllerName: t.Name(),
		},
	}
}

func (t TransactionsController) CreateTransaction(rw http.ResponseWriter, req *http.Request) {
	params := util.GetParamsFromContext(req)

	amount := 0.0
	if params["amount"] != nil {
		var ok bool
		amount, ok = params["amount"].(float64)
		if !ok {
			srverr.Handle400(rw, fmt.Errorf("invalid amount: %v", params["amount"]))
			return
		}
	}

	description := ""
	if params["description"] != nil {
		var ok bool
		description, ok = params["description"].(string)
		if !ok {
			srverr.Handle400(rw, fmt.Errorf("invalid description: %v", params["description"]))
			return
		}
	}

	groupId := ""
	if params["groupId"] != nil {
		var ok bool
		groupId, ok = params["groupId"].(string)
		if !ok {
			srverr.Handle400(rw, fmt.Errorf("invalid groupId: %v", params["groupId"]))
			return
		}
	}

	user := util.GetContextUser(req)

	newTransaction, err := domain.NewTransaction(amount, description, groupId, user.Id)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactionsModel := models.NewTransactionModel(&t.Database.Adapter)
	transaction, err := transactionsModel.Create(newTransaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transaction)
}

func (t TransactionsController) GetTransactionsByYear(rw http.ResponseWriter, req *http.Request) {
	yearString := req.PathValue("year")
	year, err := strconv.Atoi(yearString)
	if err != nil {
		srverr.Handle400(rw, fmt.Errorf("invalid year: %v", yearString))
		return
	}

	transactionsModel := models.NewTransactionModel(&t.Database.Adapter)

	user := util.GetContextUser(req)
	query := map[string]any{
		"userId": user.Id,
		"year":   year,
	}

	results, err := transactionsModel.Where(query)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactions, ok := results.([]domain.Transaction)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("results are not a slice of transactions"))
		return
	}

	transactions = sortTransactionsByTimeDesc(transactions)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transactions)
}

func (t TransactionsController) UpdateTransaction(rw http.ResponseWriter, req *http.Request) {
	params := util.GetParamsFromContext(req)

	transactionId := req.PathValue("id")

	transactionsModel := models.NewTransactionModel(&t.Database.Adapter)
	result, err := transactionsModel.Find(transactionId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transaction, ok := result.(domain.Transaction)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("result is not a transaction"))
		return
	}

	user := util.GetContextUser(req)
	if transaction.UserId.Hex() != user.Id.Hex() {
		srverr.Handle403(rw, fmt.Errorf("transaction not found"))
		return
	}

	if params["description"] != nil {
		transaction.Description = params["description"].(string)
	}
	if params["groupId"] != nil && params["groupId"] != "" {
		groupId, err := bson.ObjectIDFromHex(params["groupId"].(string))
		if err != nil {
			srverr.Handle400(rw, fmt.Errorf("invalid groupId: %v", params["groupId"]))
			return
		}
		transaction.GroupId = groupId
	}

	err = transactionsModel.Update(transaction.Id.Hex(), transaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transaction)
}

func (t TransactionsController) DeleteTransaction(rw http.ResponseWriter, req *http.Request) {
	transactionId := req.PathValue("id")

	user := util.GetContextUser(req)

	transactionsModel := models.NewTransactionModel(&t.Database.Adapter)

	result, err := transactionsModel.Find(transactionId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transaction, ok := result.(domain.Transaction)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("result is not a transaction"))
		return
	}

	if transaction.UserId.Hex() != user.Id.Hex() {
		srverr.Handle403(rw, fmt.Errorf("transaction not found"))
		return
	}

	err = transactionsModel.Delete(transactionId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func sortTransactionsByTimeDesc(transactions []domain.Transaction) []domain.Transaction {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.After(transactions[j].Timestamp)
	})

	return transactions
}

var _ Controller = TransactionsController{}
