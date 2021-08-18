package services

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserServiceName = "UserServiceName"

// UserService handles the creation, modification and deletion of users.
type UserService struct {
	MongoCN database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *UserService) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* GetAll: get all data based on condition*/
func (u *UserService) GetAll() ([]*models.User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("users")
	defer cancel()

	//create an empty array to store all fields from collection
	var users []*models.User = make([]*models.User, 0)

	//get all user record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to user variable
	for cur.Next(ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, &user)
	}
	//response data to client
	if users == nil {
		return make([]*models.User, 0), nil
	}
	return users, nil
}

/*GetOne: get one record from a collection  */
func (u *UserService) GetOne(filter bson.M) (*models.User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("users")
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	user := models.User{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &user, nil
}

/*Create: create a new record to a collection*/
func (u *UserService) Create(newUser model.NewUser) (*models.User, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("users")
	defer cancel()

	//convert newUser to bson.M
	currentTime := time.Now()
	user := models.User{
		Email:     newUser.Email,
		Password:  newUser.Password,
		Roles:     newUser.Roles,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	newData, err := utilities.InterfaceToBsonM(user)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Email:     newUser.Email,
		Roles:     newUser.Roles}, nil
}

/*UpdateOne: update one record from a collection*/
func (u UserService) UpdateOne(filter bson.M, update bson.M) (*models.User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("users")
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, errors.New("id not found")
	}

	//query the new update
	user, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return user, nil
}

//DeleteOne func is to update one record from a collection
func (u UserService) DeleteOne(filter bson.M) (*models.User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("users")
	defer cancel()

	user, errGet := u.GetOne(filter)
	if errGet != nil {
		return nil, errGet
	}

	//delete user from database
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		//response to client if there is an error.
		return nil, err
	}

	if deleteResult.DeletedCount == 0 {
		return nil, errors.New("id not found")
	}

	return user, nil
}

func (u UserService) Login(input model.Login) (*models.User, error) {
	user, err := u.GetOne(bson.M{"email": input.Email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid user or password")
	}
	//check password
	if ok := utilities.CheckPasswordHash(input.Password, user.Password); !ok {
		return nil, errors.New("invalid user or password")
	}
	return user, nil
}

//validation
func (u *UserService) ValidateNewUser(newUser model.NewUser) error {
	return validation.ValidateStruct(&newUser,
		validation.Field(&newUser.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			user, err := u.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			if user != nil {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&newUser.Password, validation.Required.Error("password must not be blanked")),
		validation.Field(&newUser.ConfirmPassword, validation.Required.Error("confirm password must not be blanked"), validation.In(newUser.Password).Error("confirm password must be identical with Password")),
		validation.Field(&newUser.Roles, validation.Required.Error("Role must not be blanked")),
	)
}

func (u *UserService) ValidateUpdateUser(id string, updateUser model.UpdateUser) error {
	return validation.ValidateStruct(&updateUser,
		validation.Field(&updateUser.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			//get current user
			currentUser, err := u.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			//check email existed or not
			user, err := u.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			if user != nil && user.Email != currentUser.Email {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&updateUser.Password, validation.Required.Error("password must not be blanked")),
		validation.Field(&updateUser.Roles, validation.Required.Error("Role must not be blanked")),
	)
}
func (u *UserService) ValidateLogin(login model.Login) error {
	return validation.ValidateStruct(&login,
		validation.Field(&login.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email")),
		validation.Field(&login.Password, validation.Required.Error("password must not be blanked")),
	)
}

func (u *UserService) HashPassword(newUser *model.NewUser) error {
	hashPassword, err := utilities.HashPassword(newUser.Password)
	if err != nil {
		return err
	}
	newUser.Password = hashPassword
	return nil
}
func (u *UserService) HashPasswordUpdateUser(updateUser *model.UpdateUser) error {
	hashPassword, err := utilities.HashPassword(updateUser.Password)
	if err != nil {
		return err
	}
	updateUser.Password = hashPassword
	return nil
}
