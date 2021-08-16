package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/graph/generated"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUser) (*models.User, error) {
	//check input
	if err := input.Validate(); err != nil {
		return nil, err
	}
	//hash password
	if err := input.HashPassword(); err != nil {
		return nil, err
	}
	newUser, err := r.userModel.Create(input)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *mutationResolver) Login(ctx context.Context, input models.Login) (*models.User, error) {
	//check input
	if err := input.Validate(); err != nil {
		return nil, err
	}
	user, err := r.userModel.Login(input)
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
	return user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users, err := r.userModel.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *queryResolver) CheckLoginStatus(ctx context.Context) (*models.User, error) {
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
	//get user based specific id
	user, err := r.userModel.GetOne(bson.M{"_id": string(id)})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
