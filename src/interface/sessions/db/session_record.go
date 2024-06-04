package sessionsdb

import (
	"golang-web-core/src/domain"
	"time"
)

type SessionRecord struct {
	Id         string `bson:"id"`
	UserId     string `bson:"userId"`
	ExpiresAt  int64  `bson:"expiresAt"`
	VerifiedIp string `bson:"verifiedIp"`
	DoesExpire bool   `bson:"doesExpire"`
}

func (r SessionRecord) IsExpired() bool {
	return r.DoesExpire && r.ExpiresAt < time.Now().UnixMilli()
}

func SessionRecordFromDomain(session domain.Session) SessionRecord {
	return SessionRecord{
		Id:         session.Id,
		UserId:     session.User.Id,
		ExpiresAt:  session.ExpiresAt,
		VerifiedIp: session.VerifiedIp,
		DoesExpire: session.DoesExpire,
	}
}
