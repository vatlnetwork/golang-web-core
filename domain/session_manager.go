package domain

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

const AuthHeaderKey string = "Authorization"
const SessionCookieName string = "session"

type ContextKey string

const ContextUserKey ContextKey = "current_user"

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

func (s SessionManager) GetContextUser(req *http.Request) (*User, error) {
	userVal := req.Context().Value(ContextUserKey)
	if userVal == nil {
		return nil, nil
	}

	user, ok := userVal.(User)
	if !ok {
		return nil, errors.New("invalid user in context")
	}

	return &user, nil
}

func (s SessionManager) SetContextUser(req *http.Request, user User) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), ContextUserKey, user))
}

func (s SessionManager) GetCurrentSession(req *http.Request) (*Session, User, error) {
	sessionId := req.Header.Get(AuthHeaderKey)
	if sessionId == "" {
		sessionCookie, err := req.Cookie(SessionCookieName)
		if err == nil {
			sessionId = sessionCookie.Value
		}
	}
	if sessionId == "" {
		return nil, User{}, nil
	}

	session, err := s.sessionRepository.GetSession(sessionId)
	if err != nil {
		if err.Error() == ErrorSessionNotFound {
			return nil, User{}, nil
		}
		return nil, User{}, err
	}

	remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, User{}, err
	}

	if !session.Validate(remoteIP) {
		return nil, User{}, nil
	}

	user, err := s.userRepository.GetUser(session.UserId)
	if err != nil {
		return nil, User{}, err
	}

	return &session, user, nil
}

func (s SessionManager) HandleSignIn(req *http.Request, email, firstName, lastName, password string, staySignedIn bool) (Session, User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	notFound := false
	if err != nil {
		if err.Error() == ErrorUserNotFound {
			notFound = true
		} else {
			return Session{}, User{}, err
		}
	}

	if notFound {
		newUser, err := NewUser(email, firstName, lastName, password)
		if err != nil {
			return Session{}, User{}, err
		}

		user, err = s.userRepository.CreateUser(newUser)
		if err != nil {
			return Session{}, User{}, err
		}
	} else {
		err = user.CheckPassword(password)
		if err != nil {
			return Session{}, User{}, err
		}
	}

	sessions, err := s.sessionRepository.GetAllForUser(user.Id)
	if err != nil {
		return Session{}, User{}, err
	}

	remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return Session{}, User{}, err
	}

	for _, session := range sessions {
		if session.Validate(remoteIP) {
			user.LastSignIn = time.Now()
			user, err = s.userRepository.UpdateUser(user)
			if err != nil {
				return Session{}, User{}, err
			}
			return session, user, nil
		}
	}

	newSession, err := NewSession(user.Id, remoteIP)
	if err != nil {
		return Session{}, User{}, err
	}

	if staySignedIn {
		newSession.Expires = false
	}

	session, err := s.sessionRepository.CreateSession(newSession)
	if err != nil {
		return Session{}, User{}, err
	}

	user.LastSignIn = time.Now()
	user, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return Session{}, User{}, err
	}

	return session, user, nil
}

func (s SessionManager) HandleSignOut(sessionId string) error {
	err := s.sessionRepository.DeleteSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}
