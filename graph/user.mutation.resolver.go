package graph

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateNewUser(input); err != nil {
		return nil, err
	}
	//hash password
	if err := service.HashPassword(&input); err != nil {
		return nil, err
	}
	newUser, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, _ := mapUser(newUser)
	return results, nil
}
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateUpdateUser(id, input); err != nil {
		return nil, err
	}
	//hash password
	if err := service.HashPasswordUpdateUser(&input); err != nil {
		return nil, err
	}
	//cast UpdateUser to bson.M type
	newUpdate, err := utilities.InterfaceToBsonM(input)
	if err != nil {
		return nil, err
	}
	updatedUser, err := service.UpdateOne(bson.M{"_id": id}, newUpdate)
	if err != nil {
		return nil, err
	}
	results, _ := mapUser(updatedUser)
	return results, nil
}
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	deletedUser, err := service.DeleteOne(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	results, _ := mapUser(deletedUser)
	return results, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateLogin(input); err != nil {
		return nil, err
	}
	user, err := service.Login(input)
	if err != nil {
		return nil, err
	}

	//encrypted for cookie value
	encryptedCookie, err := utilities.Encrypt([]byte(user.ID.Hex()))
	if err != nil {
		return nil, err
	}
	//set token
	ginContext := ctx.Value("gincontext").(*gin.Context)
	ginContext.SetCookie("netevent", string(encryptedCookie), 3600*24, "/", "localhost", false, true)
	results, _ := mapUser(user)
	return results, nil
}
func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	//set token
	ginContext := ctx.Value("gincontext").(*gin.Context)
	ginContext.SetCookie("netevent", "", 0, "/", "localhost", false, true)
	return "Logout successful", nil
}
