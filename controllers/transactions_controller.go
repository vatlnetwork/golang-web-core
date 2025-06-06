package controllers

import (
	"encoding/json"
	"errors"
	"inventory-app/domain"
	"inventory-app/srv/route"
	"inventory-app/srv/srverr"
	"inventory-app/util"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type TransactionsController struct {
	transactionRepo      domain.TransactionRepository
	sessionManager       domain.SessionManager
	transactionGroupRepo domain.TransactionGroupRepository
	moneyLocationRepo    domain.MoneyLocationRepository
}

func NewTransactionsController(
	transactionRepo domain.TransactionRepository,
	sessionManager domain.SessionManager,
	transactionGroupRepo domain.TransactionGroupRepository,
	moneyLocationRepo domain.MoneyLocationRepository,
) TransactionsController {
	return TransactionsController{
		transactionRepo:      transactionRepo,
		sessionManager:       sessionManager,
		transactionGroupRepo: transactionGroupRepo,
		moneyLocationRepo:    moneyLocationRepo,
	}
}

// BeforeAction implements Controller.
func (t TransactionsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		currentUser, err := t.sessionManager.GetContextUser(req)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		if currentUser == nil {
			srverr.Handle401(rw, errors.New("unauthorized"))
			return
		}

		handler(rw, req)
	}
}

// Name implements Controller.
func (t TransactionsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/api/transactions",
			Method:         http.MethodGet,
			Handler:        t.Index,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/location/{id}",
			Method:         http.MethodGet,
			Handler:        t.IndexByLocation,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/group/{id}",
			Method:         http.MethodGet,
			Handler:        t.IndexByGroup,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/year/{year}",
			Method:         http.MethodGet,
			Handler:        t.IndexByYear,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions",
			Method:         http.MethodPost,
			Handler:        t.Create,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/{id}",
			Method:         http.MethodGet,
			Handler:        t.Show,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/{id}",
			Method:         http.MethodPut,
			Handler:        t.Update,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/{id}",
			Method:         http.MethodDelete,
			Handler:        t.Destroy,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/upload",
			Method:         http.MethodPost,
			Handler:        t.Upload,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transactions/bulk_upload",
			Method:         http.MethodPost,
			Handler:        t.BulkUpload,
			ControllerName: t.Name(),
		},
	}
}

func (t TransactionsController) Index(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactions, err := t.transactionRepo.GetTransactionsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) IndexByLocation(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	locationId := req.PathValue("id")
	location, err := t.moneyLocationRepo.GetMoneyLocation(locationId)
	if err != nil {
		if err.Error() == domain.ErrorMoneyLocationNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if location.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorMoneyLocationNotFound))
		return
	}

	transactions, err := t.transactionRepo.GetTransactionsByLocation(locationId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) IndexByGroup(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	groupId := req.PathValue("id")
	group, err := t.transactionGroupRepo.GetTransactionGroup(groupId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionGroupNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if group.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorTransactionGroupNotFound))
		return
	}

	transactions, err := t.transactionRepo.GetTransactionsByGroup(groupId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) IndexByYear(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	year := req.PathValue("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactions, err := t.transactionRepo.GetTransactionsByYear(currentUser.Id, yearInt)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type createTransactionRequest struct {
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	GroupId         string  `json:"groupId"`
	MoneyLocationId string  `json:"moneyLocationId"`
}

func (t TransactionsController) Create(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var request createTransactionRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	newTransaction, err := domain.NewTransaction(
		currentUser.Id,
		request.Amount,
		request.Description,
		request.GroupId,
		request.MoneyLocationId,
	)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	createdTransaction, err := t.transactionRepo.CreateTransaction(newTransaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(createdTransaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) Show(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionId := req.PathValue("id")
	transaction, err := t.transactionRepo.GetTransaction(transactionId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transaction.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorTransactionNotFound))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type updateTransactionRequest struct {
	Description     string `json:"description"`
	GroupId         string `json:"groupId"`
	MoneyLocationId string `json:"moneyLocationId"`
}

func (t TransactionsController) Update(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionId := req.PathValue("id")
	transaction, err := t.transactionRepo.GetTransaction(transactionId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transaction.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorTransactionNotFound))
		return
	}

	var request updateTransactionRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transaction.Description = request.Description
	transaction.GroupId = request.GroupId
	transaction.MoneyLocationId = request.MoneyLocationId

	err = t.transactionRepo.UpdateTransaction(transaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) Destroy(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionId := req.PathValue("id")
	transaction, err := t.transactionRepo.GetTransaction(transactionId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transaction.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorTransactionNotFound))
		return
	}

	err = t.transactionRepo.DeleteTransaction(transactionId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

type uploadTransactionResponse struct {
	OldId string `json:"oldId"`
	NewId string `json:"newId"`
}

func (t TransactionsController) Upload(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var transaction domain.Transaction
	err = util.DecodeContextParams(req, &transaction)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	oldId := transaction.Id

	transaction.Id = ""
	transaction.UserId = currentUser.Id
	transaction.Year = time.UnixMilli(transaction.Timestamp).Year()

	transaction, err = t.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	response := uploadTransactionResponse{
		OldId: oldId,
		NewId: transaction.Id,
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionsController) BulkUpload(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var transactions []domain.Transaction
	err = util.DecodeContextParams(req, &transactions)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	response := []uploadTransactionResponse{}

	for _, transaction := range transactions {
		oldId := transaction.Id

		transaction.Id = ""
		transaction.UserId = currentUser.Id
		transaction.Year = time.UnixMilli(transaction.Timestamp).Year()

		transaction, err = t.transactionRepo.CreateTransaction(transaction)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		response = append(response, uploadTransactionResponse{
			OldId: oldId,
			NewId: transaction.Id,
		})
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

var _ Controller = TransactionsController{}
