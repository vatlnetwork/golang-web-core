package domain

type SessionRepository interface {
	CreateSession(session Session) (Session, error)
	GetSession(sessionId string) (Session, error)
	GetAllForUser(userId string) ([]Session, error)
	DeleteSession(sessionId string) error
	DeleteAllForUser(userId string) error
}
