package sessionrepo

import (
	"errors"
	"golang-web-core/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoSession struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    string        `bson:"userId"`
	ExpiresAt time.Time     `bson:"expiresAt"`
	Expires   bool          `bson:"expires"`
	RemoteIP  string        `bson:"remoteIP"`
}

func (s MongoSession) ToDomain() domain.Session {
	return domain.Session{
		Id:        s.Id.Hex(),
		UserId:    s.UserId,
		ExpiresAt: s.ExpiresAt,
		Expires:   s.Expires,
		RemoteIP:  s.RemoteIP,
	}
}

func MongoSessionFromDomain(session domain.Session) (MongoSession, error) {
	mongoSession := MongoSession{
		UserId:    session.UserId,
		ExpiresAt: session.ExpiresAt,
		Expires:   session.Expires,
		RemoteIP:  session.RemoteIP,
	}

	if session.Id != "" {
		id, err := bson.ObjectIDFromHex(session.Id)
		if err != nil {
			return MongoSession{}, err
		}

		if id.IsZero() {
			return MongoSession{}, errors.New("invalid session id")
		}

		mongoSession.Id = id
	}

	return mongoSession, nil
}
