package srv

import (
	media_user_int "golang-web-core/src/interface/media/user"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	users_user_int "golang-web-core/src/interface/users/user"
)

type controllers struct {
	Users    users_user_int.Controller
	Sessions sessions_user_int.Controller
	Media    media_user_int.Controller
}

func (s Server) Controllers() controllers {
	return controllers{
		Users:    users_user_int.NewUsersController(s.Repos().Users, s.Repos().Sessions, s.Repos().Media, s.Repos().OsMedia),
		Sessions: sessions_user_int.NewSessionController(s.Repos().Sessions),
		Media:    media_user_int.NewMediaController(s.Repos().Media, s.Repos().OsMedia),
	}
}
