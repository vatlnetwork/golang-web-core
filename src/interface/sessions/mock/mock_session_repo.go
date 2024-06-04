package sessionsmock

import (
	"fmt"
	"golang-web-core/src/domain"

	"github.com/google/uuid"
)

type MockSessionRepo struct{}

func NewMockSessionRepo() MockSessionRepo {
	return MockSessionRepo{}
}

var sessions []domain.Session = []domain.Session{}

func (r MockSessionRepo) FindOrCreate(session domain.Session) (domain.Session, error) {
	for _, s := range sessions {
		if s.User.Id == session.User.Id && s.VerifiedIp == session.VerifiedIp && !s.IsExpired() {
			return s, nil
		}
	}
	session.Id = uuid.NewString()
	sessions = append(sessions, session)
	return session, nil
}

func (r MockSessionRepo) Find(id string) (domain.Session, error) {
	for _, s := range sessions {
		if s.Id == id {
			return s, nil
		}
	}
	return domain.Session{}, fmt.Errorf("unable to find session with id: %v", id)
}

func (r MockSessionRepo) QueryByUserId(userId string) ([]domain.Session, error) {
	res := []domain.Session{}
	for _, s := range sessions {
		if s.User.Id == userId {
			res = append(res, s)
		}
	}
	return res, nil
}

func (r MockSessionRepo) Update(session domain.Session) error {
	index := -1
	for i := 0; i < len(sessions); i++ {
		if sessions[i].Id == session.Id {
			index = i
		}
	}
	if index < 0 || index > len(sessions)-1 {
		return fmt.Errorf("unable to find session with id: %v", session.Id)
	}
	sessions[index] = session
	return nil
}

func (r MockSessionRepo) Delete(id string) error {
	res := []domain.Session{}
	for _, s := range sessions {
		if s.Id != id {
			res = append(res, s)
		}
	}
	sessions = res
	return nil
}

func (r MockSessionRepo) DeleteExpired(userId string) error {
	sessionIds := []string{}
	for _, s := range sessions {
		if s.User.Id == userId && s.IsExpired() {
			sessionIds = append(sessionIds, s.Id)
		}
	}
	for _, id := range sessionIds {
		r.Delete(id)
	}
	return nil
}

func (r MockSessionRepo) DeleteAll(userId string) error {
	sessionIds := []string{}
	for _, s := range sessions {
		if s.User.Id == userId {
			sessionIds = append(sessionIds, s.Id)
		}
	}
	for _, id := range sessionIds {
		r.Delete(id)
	}
	return nil
}
