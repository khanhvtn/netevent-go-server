package graph

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	users, err := service.GetAll(bson.M{})
	if err != nil {
		return nil, err
	}
	results := make([]*model.User, 0)
	for _, user := range users {
		mappedUser, err := r.mapUser(user)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedUser)
	}
	return results, nil
}

func (r *queryResolver) CheckLoginStatus(ctx context.Context) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//get gin context
	ginContext := ctx.Value("gincontext").(*gin.Context)
	encryptedCookie, err := ginContext.Cookie("netevent")
	if err != nil {
		return nil, errors.New("access denied")
	}
	//decrypt cookie
	id, err := utilities.Decrypted([]byte(encryptedCookie))
	if err != nil {
		return nil, err
	}
	objectId, err := utilities.ConvertStringIdToObjectID(string(id))
	if err != nil {
		return nil, err
	}
	//get user based specific id
	user, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapUser(user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get user based specific id
	user, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapUser(user)
	if err != nil {
		return nil, err
	}
	return result, nil
}
