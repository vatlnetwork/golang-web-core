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
)

type SyncController struct {
	sessionManager       domain.SessionManager
	transactionRepo      domain.TransactionRepository
	moneyLocationRepo    domain.MoneyLocationRepository
	transactionGroupRepo domain.TransactionGroupRepository
}

func NewSyncController(
	sessionManager domain.SessionManager,
	transactionRepo domain.TransactionRepository,
	moneyLocationRepo domain.MoneyLocationRepository,
	transactionGroupRepo domain.TransactionGroupRepository,
) SyncController {
	return SyncController{
		sessionManager:       sessionManager,
		transactionRepo:      transactionRepo,
		moneyLocationRepo:    moneyLocationRepo,
		transactionGroupRepo: transactionGroupRepo,
	}
}

// BeforeAction implements Controller.
func (s SyncController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		currentUser, err := s.sessionManager.GetContextUser(req)
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
func (s SyncController) Name() string {
	return reflect.TypeOf(s).Name()
}

func (s SyncController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/api/sync/money",
			Method:         http.MethodPost,
			Handler:        s.SyncMoney,
			ControllerName: s.Name(),
		},
	}
}

type syncRequest struct {
	Transactions      []domain.Transaction      `json:"transactions"`
	MoneyLocations    []domain.MoneyLocation    `json:"moneyLocations"`
	TransactionGroups []domain.TransactionGroup `json:"transactionGroups"`
}

type syncResponse struct {
	Transactions      []domain.Transaction      `json:"transactions"`
	MoneyLocations    []domain.MoneyLocation    `json:"moneyLocations"`
	TransactionGroups []domain.TransactionGroup `json:"transactionGroups"`
}

func (s SyncController) SyncMoney(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := s.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var request syncRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	err = s.moneyLocationRepo.DeleteAllMoneyLocationsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	err = s.transactionGroupRepo.DeleteAllTransactionGroupsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	err = s.transactionRepo.DeleteAllTransactionsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	moneyLocations := request.MoneyLocations
	tranasctionGroups := request.TransactionGroups
	transactions := request.Transactions

	for i, moneyLocation := range moneyLocations {
		oldId := moneyLocation.Id

		moneyLocation.Id = ""
		moneyLocation.UserId = currentUser.Id

		moneyLocation, err = s.moneyLocationRepo.CreateMoneyLocation(moneyLocation)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		for j, transaction := range transactions {
			if transaction.MoneyLocationId == oldId {
				transactions[j].MoneyLocationId = moneyLocation.Id
			}
		}

		moneyLocations[i] = moneyLocation
	}

	for i, transactionGroup := range tranasctionGroups {
		oldId := transactionGroup.Id

		transactionGroup.Id = ""
		transactionGroup.UserId = currentUser.Id

		transactionGroup, err = s.transactionGroupRepo.CreateTransactionGroup(transactionGroup)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		for j, transaction := range transactions {
			if transaction.GroupId == oldId {
				transactions[j].GroupId = transactionGroup.Id
			}
		}

		tranasctionGroups[i] = transactionGroup
	}

	for i, transaction := range transactions {
		transaction.Id = ""
		transaction.UserId = currentUser.Id

		transaction, err = s.transactionRepo.CreateTransaction(transaction)
		if err != nil {
			srverr.Handle500(rw, err)
			return
		}

		transactions[i] = transaction
	}

	response := syncResponse{
		Transactions:      transactions,
		MoneyLocations:    moneyLocations,
		TransactionGroups: tranasctionGroups,
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

var _ Controller = SyncController{}
