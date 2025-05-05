package domain

const ErrorSessionNotFound string = "session not found"

type SessionRepository interface {
	CreateSession(session Session) (Session, error)
	GetSession(sessionId string) (Session, error)
	GetAllForUser(userId string) ([]Session, error)
	DeleteSession(sessionId string) error
	DeleteAllForUser(userId string) error
}
