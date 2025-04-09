package models

import (
	"fmt"
	"golang-web-core/app/domain"
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/mongo"
	"golang-web-core/srv/srverr"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionModel struct {
	adapter databaseadapters.DatabaseAdapter
}

func NewSessionModel(adapter databaseadapters.DatabaseAdapter) SessionModel {
	return SessionModel{
		adapter: adapter,
	}
}

// Adapter implements Model.
func (s SessionModel) Adapter() *databaseadapters.DatabaseAdapter {
	return &s.adapter
}

// All implements Model.
func (s SessionModel) All() (any, error) {
	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		cursor, err := mongoAdapter.Query(client, ctx, s.Name(), bson.M{}, nil)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		sessions := []domain.Session{}
		err = cursor.All(ctx, &sessions)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		return sessions, nil
	}

	return nil, ErrUnsupportedAdapter(s, &s.adapter)
}

// Create implements Model.
func (s SessionModel) Create(object any) (any, error) {
	session, isSession := object.(domain.Session)
	isSessionPtr := false
	if !isSession {
		var sessionPtr *domain.Session
		sessionPtr, isSessionPtr = object.(*domain.Session)
		if isSessionPtr && sessionPtr != nil {
			session = *sessionPtr
			isSession = true
		} else {
			return nil, srverr.New("the given object is not a Session or *Session")
		}
	}

	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		res, err := mongoAdapter.InsertOne(client, ctx, s.Name(), object)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		session.Id = res.InsertedID.(bson.ObjectID)
		return session, nil // Return value type
	}

	return nil, ErrUnsupportedAdapter(s, &s.adapter)
}

// Delete implements Model.
func (s SessionModel) Delete(key any) error {
	keyStr, isString := key.(string)
	if !isString {
		return srverr.New("key must be a string (ObjectID hex)")
	}

	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(keyStr)
		if err != nil {
			return srverr.Wrap(err)
		}

		filter := bson.M{
			s.PrimaryKey(): objectId,
		}

		err = mongoAdapter.DeleteOne(client, ctx, s.Name(), filter)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(s, &s.adapter)
}

// DeleteWhere implements Model.
func (s SessionModel) DeleteWhere(query map[string]any) error {
	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertSessionQueryIDs(query)
		if err != nil {
			return srverr.Wrap(err)
		}

		err = mongoAdapter.DeleteMany(client, ctx, s.Name(), filter)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(s, &s.adapter)
}

// Find implements Model.
func (s SessionModel) Find(key any) (any, error) {
	keyStr, isString := key.(string)
	if !isString {
		// Allow finding by token as well
		if token, isToken := key.(string); isToken {
			// If it's a string, try finding by token first
			foundByToken, err := s.FindByToken(token)
			if err == nil && foundByToken != nil {
				return *foundByToken, nil // Return the found session
			}
			// If not found by token or error occurred, proceed to check if it's an ID
			// Fall through to ObjectID check below
		} else {
			return nil, srverr.New("key must be a string (ObjectID hex or session token)")
		}
	}

	// If keyStr was set, attempt ObjectID conversion
	if keyStr != "" {
		mongoAdapter, ok := s.adapter.(mongo.Mongo)
		if ok {
			client, ctx, cancel, err := mongoAdapter.Connect()
			if err != nil {
				return nil, srverr.Wrap(err)
			}
			defer mongoAdapter.Close(client, ctx, cancel)

			objectId, err := bson.ObjectIDFromHex(keyStr)
			if err == nil { // Only proceed if conversion is successful
				filter := bson.M{s.PrimaryKey(): objectId}
				cursor, err := mongoAdapter.Query(client, ctx, s.Name(), filter, nil)
				if err != nil {
					return nil, srverr.Wrap(err)
				}
				sessions := []domain.Session{}
				err = cursor.All(ctx, &sessions)
				if err != nil {
					return nil, srverr.Wrap(err)
				}
				if len(sessions) > 0 {
					return sessions[0], nil
				}
			}
			// If ObjectID conversion failed OR no session found by ID, report specific error
			return nil, srverr.New(fmt.Sprintf("unable to find a Session with %s = %v or token = %v", s.PrimaryKey(), keyStr, keyStr), http.StatusNotFound)

		}
		return nil, ErrUnsupportedAdapter(s, &s.adapter)
	}

	// Should not be reached if logic is correct, but return general error
	return nil, srverr.New(fmt.Sprintf("unable to find session with key: %v", key), http.StatusNotFound)
}

// FindByToken specifically finds a session by its token string.
func (s SessionModel) FindByToken(token string) (*domain.Session, error) {
	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if !ok {
		return nil, ErrUnsupportedAdapter(s, &s.adapter)
	}

	client, ctx, cancel, err := mongoAdapter.Connect()
	if err != nil {
		return nil, srverr.Wrap(err)
	}
	defer mongoAdapter.Close(client, ctx, cancel)

	filter := bson.M{"token": token}
	cursor, err := mongoAdapter.Query(client, ctx, s.Name(), filter, nil)
	if err != nil {
		return nil, srverr.Wrap(err)
	}

	sessions := []domain.Session{}
	err = cursor.All(ctx, &sessions)
	if err != nil {
		return nil, srverr.Wrap(err)
	}

	if len(sessions) == 0 {
		return nil, srverr.New(fmt.Sprintf("no session found with token: %s", token), http.StatusNotFound) // Specific error
	}
	if len(sessions) > 1 {
		// This shouldn't happen if tokens are unique, but good to log
		fmt.Printf("Warning: Multiple sessions found for token %s\n", token)
	}

	return &sessions[0], nil
}

// Name implements Model.
func (s SessionModel) Name() string {
	return "sessions"
}

// PrimaryKey implements Model.
func (s SessionModel) PrimaryKey() string {
	return "_id"
}

// Update implements Model.
func (s SessionModel) Update(key any, object any) error {
	keyStr, isString := key.(string)
	if !isString {
		return srverr.New("key must be a string (ObjectID hex)")
	}

	isValidSessionType := false
	if _, isSession := object.(domain.Session); isSession {
		isValidSessionType = true
	} else if _, isSessionPtr := object.(*domain.Session); isSessionPtr {
		isValidSessionType = true
	}
	if !isValidSessionType {
		return srverr.New("the given object is not a Session or *Session")
	}

	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(keyStr)
		if err != nil {
			return srverr.Wrap(err)
		}

		filter := bson.M{
			s.PrimaryKey(): objectId,
		}

		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateOne(client, ctx, s.Name(), filter, update)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(s, &s.adapter)
}

// UpdateWhere implements Model.
func (s SessionModel) UpdateWhere(query map[string]any, object any) error {
	isValidSessionType := false
	if _, isSession := object.(domain.Session); isSession {
		isValidSessionType = true
	} else if _, isSessionPtr := object.(*domain.Session); isSessionPtr {
		isValidSessionType = true
	}
	if !isValidSessionType {
		return srverr.New("the given object is not a Session or *Session")
	}

	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertSessionQueryIDs(query)
		if err != nil {
			return srverr.Wrap(err)
		}

		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateMany(client, ctx, s.Name(), filter, update)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(s, &s.adapter)
}

// Where implements Model.
func (s SessionModel) Where(query map[string]any) (any, error) {
	mongoAdapter, ok := s.adapter.(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertSessionQueryIDs(query)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		cursor, err := mongoAdapter.Query(client, ctx, s.Name(), filter, nil)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		sessions := []domain.Session{}
		err = cursor.All(ctx, &sessions)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		return sessions, nil
	}
	return nil, ErrUnsupportedAdapter(s, &s.adapter)
}

// Helper function specific to SessionModel queries involving ObjectIDs
func convertSessionQueryIDs(query map[string]any) (bson.M, error) {
	filter := bson.M{}
	for key, value := range query {
		// Convert known ObjectID fields if they are strings
		if key == "_id" || key == "userId" { // Add other ObjectID fields here if needed
			if idStr, ok := value.(string); ok {
				objID, err := bson.ObjectIDFromHex(idStr)
				if err != nil {
					return nil, fmt.Errorf("invalid ObjectID format for key '%s': %v", key, err)
				}
				filter[key] = objID
			} else {
				// Pass through if not a string (might be already ObjectID)
				filter[key] = value
			}
		} else {
			// Pass through other fields directly (like 'token', 'expiresAt')
			filter[key] = value
		}
	}
	return filter, nil
}

var _ Model = SessionModel{} // Compile-time check
