package util

import (
	"fmt"
	"golang-web-core/app/domain"
	"net/http"
)

func GetContextUser(req *http.Request) *domain.User {
	user, ok := req.Context().Value("current_user").(domain.User)
	if !ok {
		return nil
	}

	return &user
}

func ContextUserOrError(req *http.Request) (domain.User, error) {
	user := GetContextUser(req)
	if user == nil {
		return domain.User{}, fmt.Errorf("not logged in")
	}

	return *user, nil
}
