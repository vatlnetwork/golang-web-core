package users_user_int

import (
	"encoding/json"
	"golang-web-core/src/domain"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/srverr"
	"io"
	"net/http"
)

type UpdateThemeForm struct {
	Theme      domain.Theme `json:"theme"`
	ThemeColor domain.Color `json:"themeColor"`
}

func (c Controller) UpdateTheme(rw http.ResponseWriter, req *http.Request) {
	// decode request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}
	formData := UpdateThemeForm{}
	err = json.Unmarshal(body, &formData)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusBadRequest)
		return
	}

	// get current user
	currentUser, err := sessions_user_int.CurrentUser(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// update current user theme
	currentUser.Theme = formData.Theme
	currentUser.ThemeColor = formData.ThemeColor
	err = c.UserRepository.Update(currentUser)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
