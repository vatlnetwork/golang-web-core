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

type UserModel struct {
	adapter *databaseadapters.DatabaseAdapter
}

func NewUserModel(adapter *databaseadapters.DatabaseAdapter) UserModel {
	return UserModel{
		adapter: adapter,
	}
}

// Adapter implements Model.
func (u UserModel) Adapter() *databaseadapters.DatabaseAdapter {
	return u.adapter
}

// All implements Model.
func (u UserModel) All() (any, error) {
	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		cursor, err := mongoAdapter.Query(client, ctx, u.Name(), bson.M{}, nil)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		users := []domain.User{}
		err = cursor.All(ctx, &users)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		return users, nil
	}

	return nil, ErrUnsupportedAdapter(u, u.adapter)
}

// Create implements Model.
func (u UserModel) Create(object any) (any, error) {
	user, isUser := object.(domain.User)
	if !isUser {
		// Attempt pointer conversion if direct type assertion fails
		if userPtr, isUserPtr := object.(*domain.User); isUserPtr && userPtr != nil {
			user = *userPtr
			isUser = true
		} else {
			return nil, srverr.New("the given object is not a User or *User")
		}
	}

	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		// Pass the original object (which might be *User) or the dereferenced user
		res, err := mongoAdapter.InsertOne(client, ctx, u.Name(), object)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		user.Id = res.InsertedID.(bson.ObjectID)
		return user, nil // Return the value type domain.User
	}

	return nil, ErrUnsupportedAdapter(u, u.adapter)
}

// Delete implements Model.
func (u UserModel) Delete(key any) error {
	keyStr, isString := key.(string)
	if !isString {
		return srverr.New("key must be a string")
	}

	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
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
			u.PrimaryKey(): objectId,
		}

		err = mongoAdapter.DeleteOne(client, ctx, u.Name(), filter)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(u, u.adapter)
}

// DeleteWhere implements Model.
func (u UserModel) DeleteWhere(query map[string]any) error {
	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertQueryIDs(query)
		if err != nil {
			return srverr.Wrap(err)
		}

		err = mongoAdapter.DeleteMany(client, ctx, u.Name(), filter)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(u, u.adapter)
}

// Find implements Model.
func (u UserModel) Find(key any) (any, error) {
	keyStr, isString := key.(string)
	if !isString {
		return nil, srverr.New("key must be a string")
	}

	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		objectId, err := bson.ObjectIDFromHex(keyStr)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		filter := bson.M{
			u.PrimaryKey(): objectId,
		}

		cursor, err := mongoAdapter.Query(client, ctx, u.Name(), filter, nil)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		users := []domain.User{}
		err = cursor.All(ctx, &users)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		if len(users) == 0 {
			return nil, srverr.New(fmt.Sprintf("unable to find a User with %v = %v", u.PrimaryKey(), key), http.StatusNotFound)
		}

		return users[0], nil
	}

	return nil, ErrUnsupportedAdapter(u, u.adapter)
}

func (u UserModel) FindByEmail(email string) (domain.User, error) {
	results, err := u.Where(map[string]any{"email": email})
	if err != nil {
		return domain.User{}, srverr.Wrap(err)
	}

	users, ok := results.([]domain.User)
	if !ok {
		return domain.User{}, srverr.New("results are not a slice of User")
	}

	if len(users) == 0 {
		return domain.User{}, srverr.New(fmt.Sprintf("unable to find a User with email %v", email), http.StatusNotFound)
	}

	return users[0], nil
}

// Name implements Model.
func (u UserModel) Name() string {
	return "users"
}

// PrimaryKey implements Model.
func (u UserModel) PrimaryKey() string {
	return "_id"
}

// Update implements Model.
func (u UserModel) Update(key any, object any) error {
	keyStr, isString := key.(string)
	if !isString {
		return srverr.New("key must be a string")
	}

	_, isUser := object.(domain.User)
	if !isUser {
		// Check for pointer type as well
		_, isUserPtr := object.(*domain.User)
		if !isUserPtr {
			return srverr.New("the given object is not a User or *User")
		}
		isUser = true // Mark as valid if it's a pointer
	}

	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
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
			u.PrimaryKey(): objectId,
		}

		// Use the passed object directly in $set, mongo driver handles value/pointer
		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateOne(client, ctx, u.Name(), filter, update)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(u, u.adapter)
}

// UpdateWhere implements Model.
func (u UserModel) UpdateWhere(query map[string]any, object any) error {
	_, isUser := object.(domain.User)
	if !isUser {
		// Check for pointer type as well
		_, isUserPtr := object.(*domain.User)
		if !isUserPtr {
			return srverr.New("the given object is not a User or *User")
		}
		isUser = true // Mark as valid if it's a pointer
	}

	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertQueryIDs(query)
		if err != nil {
			return srverr.Wrap(err)
		}

		// Use the passed object directly in $set
		update := bson.M{
			"$set": object,
		}

		err = mongoAdapter.UpdateMany(client, ctx, u.Name(), filter, update)
		if err != nil {
			return srverr.Wrap(err)
		}

		return nil
	}

	return ErrUnsupportedAdapter(u, u.adapter)
}

// Where implements Model.
func (u UserModel) Where(query map[string]any) (any, error) {
	mongoAdapter, ok := (*u.adapter).(mongo.Mongo)
	if ok {
		client, ctx, cancel, err := mongoAdapter.Connect()
		if err != nil {
			return nil, srverr.Wrap(err)
		}
		defer mongoAdapter.Close(client, ctx, cancel)

		filter, err := convertQueryIDs(query)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		cursor, err := mongoAdapter.Query(client, ctx, u.Name(), filter, nil)
		if err != nil {
			return nil, err
		}

		users := []domain.User{}
		err = cursor.All(ctx, &users)
		if err != nil {
			return nil, srverr.Wrap(err)
		}

		return users, nil
	}
	return nil, ErrUnsupportedAdapter(u, u.adapter)
}

// Helper function to convert known ID fields in a query map to bson.ObjectID
func convertQueryIDs(query map[string]any) (bson.M, error) {
	filter := bson.M{}
	for key, value := range query {
		if key == "_id" || key == "userId" { // Add other potential ID fields here
			if idStr, ok := value.(string); ok {
				objID, err := bson.ObjectIDFromHex(idStr)
				if err != nil {
					return nil, fmt.Errorf("invalid ObjectID format for key '%s': %v", key, err)
				}
				filter[key] = objID
			} else {
				// If it's not a string, pass it through (could already be ObjectID or other type)
				filter[key] = value
			}
		} else {
			filter[key] = value
		}
	}
	return filter, nil
}

var _ Model = UserModel{}
