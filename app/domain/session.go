package domain

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Session struct {
	Id        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    bson.ObjectID `json:"userId" bson:"userId"`
	Token     string        `json:"token" bson:"token"`
	ExpiresAt time.Time     `json:"expiresAt" bson:"expiresAt"`
	Expires   bool          `json:"expires" bson:"expires"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func NewSession(userId bson.ObjectID) (Session, error) {
	session := Session{
		UserId:    userId,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		Expires:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return session, nil
}

func (s Session) IsExpired() bool {
	return s.Expires && time.Now().After(s.ExpiresAt)
}

func (s *Session) ResetExpiration() {
	s.ExpiresAt = time.Now().Add(time.Hour * 24 * 30)
	s.UpdatedAt = time.Now()
}
