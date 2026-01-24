package domain

import (
	"errors"
	"time"
)

const ErrorSessionNotFound string = "session not found"

type SessionRepository interface {
	CreateSession(session Session) (Session, error)
	GetSession(sessionId string) (Session, error)
	GetAllForUser(userId string) ([]Session, error)
	DeleteSession(sessionId string) error
	DeleteAllForUser(userId string) error
}

type Session struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Expires   bool      `json:"expires"`
	RemoteIP  string    `json:"remoteIP"`
	UserAgent string    `json:"userAgent"`
}

func NewSession(userId, remoteIP, userAgent string) (Session, error) {
	if userId == "" {
		return Session{}, errors.New("user id is required")
	}

	if remoteIP == "" {
		return Session{}, errors.New("remote ip is required")
	}

	if userAgent == "" {
		return Session{}, errors.New("user agent is required")
	}

	return Session{
		UserId:    userId,
		RemoteIP:  remoteIP,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		Expires:   true,
	}, nil
}

func (s Session) IsExpired() bool {
	return s.Expires && time.Now().After(s.ExpiresAt)
}

func (s Session) Validate(remoteIP string) bool {
	return s.RemoteIP == remoteIP && !s.IsExpired()
}
