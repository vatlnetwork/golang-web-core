package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Expires   bool      `json:"expires"`
	RemoteIP  string    `json:"remoteIP"`
}

func NewSession(userId, remoteIP string) (Session, error) {
	if userId == "" {
		return Session{}, errors.New("user id is required")
	}

	if remoteIP == "" {
		return Session{}, errors.New("remote ip is required")
	}

	return Session{
		UserId:    userId,
		RemoteIP:  remoteIP,
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
