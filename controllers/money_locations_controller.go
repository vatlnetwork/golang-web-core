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

type MoneyLocationsController struct {
	moneyLocationRepo domain.MoneyLocationRepository
	sessionManager    domain.SessionManager
	transactionRepo   domain.TransactionRepository
}

func NewMoneyLocationsController(
	moneyLocationRepo domain.MoneyLocationRepository,
	sessionManager domain.SessionManager,
	transactionRepo domain.TransactionRepository,
) MoneyLocationsController {
	return MoneyLocationsController{
		moneyLocationRepo: moneyLocationRepo,
		sessionManager:    sessionManager,
		transactionRepo:   transactionRepo,
	}
}

// BeforeAction implements Controller.
func (m MoneyLocationsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		currentUser, err := m.sessionManager.GetContextUser(req)
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
func (m MoneyLocationsController) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m MoneyLocationsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/api/money_locations",
			Method:         http.MethodGet,
			Handler:        m.Index,
			ControllerName: m.Name(),
		},
		{
			Pattern:        "/api/money_locations",
			Method:         http.MethodPost,
			Handler:        m.Create,
			ControllerName: m.Name(),
		},
		{
			Pattern:        "/api/money_locations/{id}",
			Method:         http.MethodGet,
			Handler:        m.Show,
			ControllerName: m.Name(),
		},
		{
			Pattern:        "/api/money_locations/{id}",
			Method:         http.MethodPut,
			Handler:        m.Update,
			ControllerName: m.Name(),
		},
		{
			Pattern:        "/api/money_locations/{id}",
			Method:         http.MethodDelete,
			Handler:        m.Destroy,
			ControllerName: m.Name(),
		},
	}
}

func (m MoneyLocationsController) Index(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	locations, err := m.moneyLocationRepo.GetMoneyLocationsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(locations)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type moneyLocationCreateRequest struct {
	Name string `json:"name"`
}

func (m MoneyLocationsController) Create(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var request moneyLocationCreateRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	newLocation, err := domain.NewMoneyLocation(request.Name, currentUser.Id)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	location, err := m.moneyLocationRepo.CreateMoneyLocation(newLocation)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(location)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (m MoneyLocationsController) Show(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	locationId := req.PathValue("id")
	location, err := m.moneyLocationRepo.GetMoneyLocation(locationId)
	if err != nil {
		if err.Error() == domain.ErrorMoneyLocationNotFound {
			srverr.Handle404(rw, err)
			return
		}
		srverr.Handle500(rw, err)
		return
	}

	if location.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorMoneyLocationNotFound))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(location)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type updateRequest struct {
	Name string `json:"name"`
}

func (m MoneyLocationsController) Update(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	locationId := req.PathValue("id")
	location, err := m.moneyLocationRepo.GetMoneyLocation(locationId)
	if err != nil {
		if err.Error() == domain.ErrorMoneyLocationNotFound {
			srverr.Handle404(rw, err)
			return
		}
		srverr.Handle500(rw, err)
		return
	}

	if location.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorMoneyLocationNotFound))
		return
	}

	var request updateRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	location.Name = request.Name
	location, err = m.moneyLocationRepo.UpdateMoneyLocation(location)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(location)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type deleteRequest struct {
	DeleteTransactions bool `json:"deleteTransactions"`
}

func (m MoneyLocationsController) Destroy(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	locationId := req.PathValue("id")
	location, err := m.moneyLocationRepo.GetMoneyLocation(locationId)
	if err != nil {
		if err.Error() == domain.ErrorMoneyLocationNotFound {
			srverr.Handle404(rw, err)
			return
		}
		srverr.Handle500(rw, err)
		return
	}

	if location.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New(domain.ErrorMoneyLocationNotFound))
		return
	}

	var request deleteRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactions, err := m.transactionRepo.GetTransactionsByLocation(locationId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	if len(transactions) > 0 && !request.DeleteTransactions {
		srverr.Handle400(rw, errors.New("this location has transactions, set deleteTransactions to true to delete them"))
		return
	}

	err = m.transactionRepo.DeleteTransactionsInLocation(locationId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	err = m.moneyLocationRepo.DeleteMoneyLocation(locationId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

var _ Controller = MoneyLocationsController{}
