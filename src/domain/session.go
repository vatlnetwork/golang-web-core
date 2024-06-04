package domain

import (
	"fmt"
	"golang-web-core/src/util"
	"net"
	"time"

	"github.com/google/uuid"
)

type SessionKey string

const CurrentSessionKey SessionKey = "currentSession"

type Session struct {
	Id         string `json:"id" bson:"id"`
	User       User   `json:"user" bson:"user"`
	ExpiresAt  int64  `json:"expiresAt" bson:"expiresAt"`
	VerifiedIp string `json:"verifiedIp" bson:"verifiedIp"`
	DoesExpire bool   `json:"doesExpire" bson:"doesExpire"`
}

func (s Session) IsExpired() bool {
	return s.ExpiresAt < time.Now().UnixMilli() && s.DoesExpire
}

func (s Session) IsValid(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("error spliting host and port in Session.IsValid: %v\n", err.Error())
		return false
	}
	return !s.IsExpired() && host == s.VerifiedIp
}

func NewSession(u User, remoteAddr string, expires bool) (Session, error) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return Session{}, err
	}
	return Session{
		Id:         uuid.NewString(),
		User:       u,
		ExpiresAt:  time.Now().UnixMilli() + util.DaysInMilliseconds(1),
		VerifiedIp: host,
		DoesExpire: expires,
	}, nil
}
