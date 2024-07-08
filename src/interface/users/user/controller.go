package users_user_int

import (
	"golang-web-core/src/application/srv/routes"
	"golang-web-core/src/domain"
	media_os_int "golang-web-core/src/interface/media/os"
	"net/http"
)

type Controller struct {
	domain.UserRepository
	domain.SessionRepository
	domain.MediaRepository
	media_os_int.OsMediaRepo
}

func NewUsersController(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	medRepo domain.MediaRepository,
	medOsRepo media_os_int.OsMediaRepo,
) Controller {
	return Controller{
		UserRepository:    userRepo,
		SessionRepository: sessionRepo,
		MediaRepository:   medRepo,
		OsMediaRepo:       medOsRepo,
	}
}

func (c Controller) Routes() []routes.Route {
	return []routes.Route{
		{
			Path:    "/users/sign_in",
			Method:  http.MethodPost,
			Handler: c.SignIn,
		},
		{
			Path:    "/users/sign_up",
			Method:  http.MethodPost,
			Handler: c.SignUp,
		},
		{
			Path:         "/users/update_theme",
			Method:       http.MethodPost,
			Handler:      c.UpdateTheme,
			RequiresAuth: true,
		},
		{
			Path:         "/users/update",
			Method:       http.MethodPost,
			Handler:      c.Update,
			RequiresAuth: true,
		},
		{
			Path:         "/users/update_profile_picture",
			Method:       http.MethodPost,
			Handler:      c.UpdateProfilePicture,
			RequiresAuth: true,
		},
		{
			Path:         "/users/update_password",
			Method:       http.MethodPost,
			Handler:      c.UpdatePassword,
			RequiresAuth: true,
		},
		{
			Path:         "/users/update_address",
			Method:       http.MethodPost,
			Handler:      c.UpdateAddress,
			RequiresAuth: true,
		},
		{
			Path:         "/users/delete",
			Method:       http.MethodDelete,
			Handler:      c.DeleteUser,
			RequiresAuth: true,
		},
		{
			Path:         "/users/confirm_password",
			Method:       http.MethodPost,
			Handler:      c.ConfirmPassword,
			RequiresAuth: true,
		},
	}
}
