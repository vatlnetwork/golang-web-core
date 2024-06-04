package srv

import (
	"golang-web-core/src/domain"
	mediadb "golang-web-core/src/interface/media/db"
	mediamock "golang-web-core/src/interface/media/mock"
	media_os_int "golang-web-core/src/interface/media/os"
	sessionsdb "golang-web-core/src/interface/sessions/db"
	sessionsmock "golang-web-core/src/interface/sessions/mock"
	usersdb "golang-web-core/src/interface/users/db"
	usersmock "golang-web-core/src/interface/users/mock"
	"golang-web-core/src/srv/cfg"
)

type repos struct {
	Users    domain.UserRepository
	Sessions domain.SessionRepository
	Media    domain.MediaRepository
	OsMedia  media_os_int.OsMediaRepo
}

func (s Server) Repos() repos {
	switch s.Config.Env {
	case cfg.MockDev:
		return s.MockRepos()
	default:
		return s.IntegRepos()
	}
}

func (s Server) MockRepos() repos {
	return repos{
		Users:    usersmock.NewMockUserRepo(),
		Sessions: sessionsmock.NewMockSessionRepo(),
		Media:    mediamock.NewMockMediaRepo(),
		OsMedia:  media_os_int.NewOsMediaRepo(s.Config.Media),
	}
}

func (s Server) IntegRepos() repos {
	return repos{
		Users:    usersdb.NewMongoUserRepository(s.Config.Mongo),
		Sessions: sessionsdb.NewMongoSessionRepository(s.Config.Mongo),
		Media:    mediadb.NewMongoMediaRepo(s.Config.Mongo),
		OsMedia:  media_os_int.NewOsMediaRepo(s.Config.Media),
	}
}
