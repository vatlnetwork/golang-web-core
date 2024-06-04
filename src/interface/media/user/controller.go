package media_user_int

import (
	"golang-web-core/src/domain"
	media_os_int "golang-web-core/src/interface/media/os"
	"golang-web-core/src/srv/routes"
	"net/http"
)

type Controller struct {
	domain.MediaRepository
	media_os_int.OsMediaRepo
}

func NewMediaController(medRepo domain.MediaRepository, osMedRepo media_os_int.OsMediaRepo) Controller {
	return Controller{
		MediaRepository: medRepo,
		OsMediaRepo:     osMedRepo,
	}
}

func (c Controller) Routes() []routes.Route {
	return []routes.Route{
		{
			Path:         "/media/upload",
			Method:       http.MethodPost,
			Handler:      c.Upload,
			RequiresAuth: true,
		},
		{
			Path:    "/media/{id}",
			Method:  http.MethodGet,
			Handler: c.Get,
		},
		{
			Path:         "/media/delete/{id}",
			Method:       http.MethodDelete,
			Handler:      c.Delete,
			RequiresAuth: true,
		},
		{
			Path:    "/media/user/{id}",
			Method:  http.MethodGet,
			Handler: c.UserMedia,
		},
	}
}
