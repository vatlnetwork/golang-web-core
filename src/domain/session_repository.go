package domain

type SessionRepository interface {
	FindOrCreate(session Session) (Session, error)
	Find(id string) (Session, error)
	QueryByUserId(userId string) ([]Session, error)
	Update(session Session) error
	Delete(id string) error
	DeleteExpired(userId string) error
	DeleteAll(userId string) error
}
