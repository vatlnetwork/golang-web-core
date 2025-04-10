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
)

type TransactionGroupsController struct {
	cfg.Config
}

func NewTransactionGroupsController(cfg cfg.Config) TransactionGroupsController {
	return TransactionGroupsController{
		Config: cfg,
	}
}

func (t TransactionGroupsController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, err := util.ContextUserOrError(req)
		if err != nil {
			srverr.Handle401(rw, err)
			return
		}

		handler(rw, req)
	}
}

// Name implements Controller.
func (t TransactionGroupsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionGroupsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/transaction_groups",
			Method:         http.MethodPost,
			Handler:        t.CreateTransactionGroup,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transaction_groups",
			Method:         http.MethodGet,
			Handler:        t.Index,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transaction_groups/{id}",
			Method:         http.MethodPut,
			Handler:        t.UpdateTransactionGroup,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/transaction_groups/{id}",
			Method:         http.MethodDelete,
			Handler:        t.DeleteTransactionGroup,
			ControllerName: t.Name(),
		},
	}
}

func (t TransactionGroupsController) CreateTransactionGroup(rw http.ResponseWriter, req *http.Request) {
	params, err := util.GetParams(req)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	description := ""
	if params["description"] != nil {
		description = params["description"].(string)
	}

	if description == "" {
		description = "Unnamed Group"
	}

	user := util.GetContextUser(req)

	newTransactionGroup := domain.NewTransactionGroup(description, user.Id)

	transactionGroupsModel := models.NewTransactionGroupModel(&t.Database.Adapter)
	transactionGroup, err := transactionGroupsModel.Create(newTransactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transactionGroup)
}

func (t TransactionGroupsController) Index(rw http.ResponseWriter, req *http.Request) {
	transactionGroupsModel := models.NewTransactionGroupModel(&t.Database.Adapter)
	results, err := transactionGroupsModel.All()
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroups, ok := results.([]domain.TransactionGroup)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("results are not a slice of transaction groups"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transactionGroups)
}

func (t TransactionGroupsController) UpdateTransactionGroup(rw http.ResponseWriter, req *http.Request) {
	params, err := util.GetParams(req)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactionGroupId := req.PathValue("id")

	transactionGroupsModel := models.NewTransactionGroupModel(&t.Database.Adapter)
	result, err := transactionGroupsModel.Find(transactionGroupId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroup, ok := result.(domain.TransactionGroup)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("result is not a transaction group"))
		return
	}

	if params["description"] != nil {
		transactionGroup.Description = params["description"].(string)
	}

	err = transactionGroupsModel.Update(transactionGroupId, transactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(transactionGroup)
}

func (t TransactionGroupsController) DeleteTransactionGroup(rw http.ResponseWriter, req *http.Request) {
	transactionGroupId := req.PathValue("id")

	transactionGroupsModel := models.NewTransactionGroupModel(&t.Database.Adapter)
	err := transactionGroupsModel.Delete(transactionGroupId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

var _ Controller = TransactionGroupsController{}
