package sessions_user_int

import (
	"fmt"
	"golang-web-core/src/domain"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// use the current request context to get the session from the database
func FetchCurrentSession(req *http.Request, sessionsRepo domain.SessionRepository) (*domain.Session, error) {
	sessionId := req.Header.Get("Authorization")
	session, err := sessionsRepo.Find(sessionId)
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			return nil, err
		}
	}
	if session.IsValid(req.RemoteAddr) {
		return &session, nil
	}
	return nil, nil
}

const BadConversionError string = "unable to convert request context session into a domain session"

// get the current session out of the request context
// the context must contain the session, not just the id in the authorization header
func CurrentSession(req *http.Request) (*domain.Session, error) {
	currentSessionVal := req.Context().Value(domain.CurrentSessionKey)
	if currentSessionVal == nil {
		return nil, nil
	}
	currentSession, ok := currentSessionVal.(domain.Session)
	if !ok {
		return nil, fmt.Errorf(BadConversionError)
	}
	return &currentSession, nil
}

// get the current user out of the current session (also requires current session to be in context)
func CurrentUser(req *http.Request) (domain.User, error) {
	currentSession, err := CurrentSession(req)
	if err != nil {
		return domain.User{}, err
	}
	if currentSession == nil {
		return domain.User{}, fmt.Errorf("currentSession is nil")
	}
	currentUser := (*currentSession).User
	return currentUser, nil
}
