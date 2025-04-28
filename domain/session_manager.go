package domain

import "net/http"

type SessionManager struct {
	sessionRepository SessionRepository
	userRepository    UserRepository
}

func NewSessionManager(sessionRepository SessionRepository, userRepository UserRepository) SessionManager {
	return SessionManager{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (s SessionManager) GetCurrentSession(req *http.Request) (Session, User, error) {
	panic("not implemented")
}

func (s SessionManager) HandleSignIn(req *http.Request, email string, password string) (Session, User, error) {
	panic("not implemented")
}

func (s SessionManager) HandleSignOut(userId string) error {
	panic("not implemented")
}
