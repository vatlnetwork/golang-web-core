package controllers

import (
	"fmt"
	"inventory-app/domain"
	moneylocationrepo "inventory-app/repositories/money_location"
	"inventory-app/srv/cfg"
	"inventory-app/srv/route"
	"net/http"
	"reflect"
)

type MoneyLocationsController struct {
	moneyLocationRepo domain.MoneyLocationRepository
}

func NewMoneyLocationsController(moneyLocationRepo domain.MoneyLocationRepository) MoneyLocationsController {
	return MoneyLocationsController{
		moneyLocationRepo: moneyLocationRepo,
	}
}

func NewMoneyLocationsControllerFromConfig(config cfg.Config) (MoneyLocationsController, error) {
	var moneyLocationRepo domain.MoneyLocationRepository

	switch config.MoneyLocationRepository {
	case "MongoMoneyLocationRepository":
		if !config.Mongo.IsEnabled() {
			return MoneyLocationsController{}, fmt.Errorf("mongo is not enabled")
		}
		moneyLocationRepo = moneylocationrepo.NewMongoMoneyLocationRepository(config.Mongo, config.Env == cfg.Development)
	default:
		return MoneyLocationsController{}, fmt.Errorf("invalid money location repository: %v", config.MoneyLocationRepository)
	}

	return NewMoneyLocationsController(moneyLocationRepo), nil
}

// BeforeAction implements Controller.
func (m MoneyLocationsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

// Name implements Controller.
func (m MoneyLocationsController) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m MoneyLocationsController) Routes() []route.Route {
	return []route.Route{}
}

var _ Controller = MoneyLocationsController{}
