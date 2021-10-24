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
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newUser, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	//send activate account mail
	if err := utilities.SendActivateAccountMail(newUser); err != nil {
		return nil, err
	}
	results, err := r.mapUser(newUser)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) ActivateUser(ctx context.Context, input model.ActivateUser) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateActivateUser(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	//hash password
	if err := service.HashPassword(&input); err != nil {
		return nil, err
	}

	//convert userID to objectId
	objectId, err := utilities.ConvertStringIdToObjectID(input.UserID)
	if err != nil {
		return nil, err
	}
	activatedUser, err := service.ActivateAccount(*objectId,
		bson.M{
			"password":    input.Password,
			"isActivated": true,
		},
	)
	if err != nil {
		return nil, err
	}
	results, err := r.mapUser(activatedUser)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateUpdateUser(id, input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
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
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	updatedUser, err := service.UpdateOne(bson.M{"_id": objectId}, newUpdate)
	if err != nil {
		return nil, err
	}
	results, err := r.mapUser(updatedUser)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedUser, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapUser(deletedUser)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	//check input
	if err := service.ValidateLogin(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
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
	results, err := r.mapUser(user)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	//set token
	ginContext := ctx.Value("gincontext").(*gin.Context)
	ginContext.SetCookie("netevent", "", 0, "/", "localhost", false, true)
	return "Logout successful", nil
}
