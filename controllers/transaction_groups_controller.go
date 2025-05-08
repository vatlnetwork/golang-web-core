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
func (t TransactionGroupsController) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t TransactionGroupsController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/api/transaction_groups",
			Method:         http.MethodGet,
			Handler:        t.Index,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transaction_groups",
			Method:         http.MethodPost,
			Handler:        t.Create,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transaction_groups/{id}",
			Method:         http.MethodGet,
			Handler:        t.Show,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transaction_groups/{id}",
			Method:         http.MethodPut,
			Handler:        t.Update,
			ControllerName: t.Name(),
		},
		{
			Pattern:        "/api/transaction_groups/{id}",
			Method:         http.MethodDelete,
			Handler:        t.Destroy,
			ControllerName: t.Name(),
		},
	}
}

func (t TransactionGroupsController) Index(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroups, err := t.transactionGroupRepo.GetTransactionGroupsForUser(currentUser.Id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactionGroups)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type transactionGroupCreateRequest struct {
	Description string `json:"description"`
}

func (t TransactionGroupsController) Create(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	var request transactionGroupCreateRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	newTransactionGroup, err := domain.NewTransactionGroup(currentUser.Id, request.Description)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactionGroup, err := t.transactionGroupRepo.CreateTransactionGroup(newTransactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionGroupsController) Show(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroupId := req.PathValue("id")
	transactionGroup, err := t.transactionGroupRepo.GetTransactionGroup(transactionGroupId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionGroupNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transactionGroup.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New("transaction group not found"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

type transactionGroupUpdateRequest struct {
	Description string `json:"description"`
}

func (t TransactionGroupsController) Update(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroupId := req.PathValue("id")
	transactionGroup, err := t.transactionGroupRepo.GetTransactionGroup(transactionGroupId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionGroupNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transactionGroup.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New("transaction group not found"))
		return
	}

	var request transactionGroupUpdateRequest
	err = util.DecodeContextParams(req, &request)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	transactionGroup.Description = request.Description
	err = t.transactionGroupRepo.UpdateTransactionGroup(transactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(transactionGroup)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}

func (t TransactionGroupsController) Destroy(rw http.ResponseWriter, req *http.Request) {
	currentUser, err := t.sessionManager.GetContextUser(req)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	transactionGroupId := req.PathValue("id")
	transactionGroup, err := t.transactionGroupRepo.GetTransactionGroup(transactionGroupId)
	if err != nil {
		if err.Error() == domain.ErrorTransactionGroupNotFound {
			srverr.Handle404(rw, err)
		} else {
			srverr.Handle500(rw, err)
		}
		return
	}

	if transactionGroup.UserId != currentUser.Id {
		srverr.Handle404(rw, errors.New("transaction group not found"))
		return
	}

	err = t.transactionGroupRepo.DeleteTransactionGroup(transactionGroupId)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

var _ Controller = TransactionGroupsController{}
