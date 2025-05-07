package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"inventory-app/domain"
	moneylocationrepo "inventory-app/repositories/money_location"
	sessionrepo "inventory-app/repositories/session"
	userrepo "inventory-app/repositories/user"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"inventory-app/srv/srverr"
	"inventory-app/util"
	"net/http"
	"reflect"
)

type MoneyLocationsController struct {
	moneyLocationRepo domain.MoneyLocationRepository
	sessionManager    domain.SessionManager
}

func NewMoneyLocationsController(moneyLocationRepo domain.MoneyLocationRepository, sessionManager domain.SessionManager) MoneyLocationsController {
	return MoneyLocationsController{
		moneyLocationRepo: moneyLocationRepo,
		sessionManager:    sessionManager,
	}
}

func NewMoneyLocationsControllerFromConfig(config cfg.Config) (MoneyLocationsController, error) {
	var moneyLocationRepo domain.MoneyLocationRepository
	var sessionRepo domain.SessionRepository
	var userRepo domain.UserRepository

	switch config.MoneyLocationRepository {
	case "MongoMoneyLocationRepository":
		if !config.Mongo.IsEnabled() {
			return MoneyLocationsController{}, fmt.Errorf("mongo is not enabled")
		}
		moneyLocationRepo = moneylocationrepo.NewMongoMoneyLocationRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return MoneyLocationsController{}, fmt.Errorf("invalid money location repository: %v", config.MoneyLocationRepository)
	}

	switch config.SessionRepository {
	case "MongoSessionRepository":
		if !config.Mongo.IsEnabled() {
			return MoneyLocationsController{}, fmt.Errorf("mongo is not enabled")
		}
		sessionRepo = sessionrepo.NewMongoSessionRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return MoneyLocationsController{}, fmt.Errorf("invalid session repository: %v", config.SessionRepository)
	}

	switch config.UserRepository {
	case "MongoUserRepository":
		if !config.Mongo.IsEnabled() {
			return MoneyLocationsController{}, fmt.Errorf("mongo is not enabled")
		}
		userRepo = userrepo.NewMongoUserRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return MoneyLocationsController{}, fmt.Errorf("invalid user repository: %v", config.UserRepository)
	}

	sessionManager := domain.NewSessionManager(sessionRepo, userRepo)

	return NewMoneyLocationsController(moneyLocationRepo, sessionManager), nil
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

type createRequest struct {
	Name string `json:"name"`
}

func (m MoneyLocationsController) Create(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := m.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var request createRequest
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

	err = m.moneyLocationRepo.DeleteMoneyLocation(locationId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

var _ Controller = MoneyLocationsController{}
