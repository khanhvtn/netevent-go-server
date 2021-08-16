package models

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/* Model Type */
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password,omitempty"`
	Roles     []string           `bson:"roles" json:"roles"`
}

/* Input Type */
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Email           string   `json:"email"`
	Password        string   `json:"password"`
	ConfirmPassword string   `json:"confirmPassword"`
	Roles           []string `json:"roles"`
}

/* Model Function */
/* createContextAndTargetCol: create and return targeted collection based on collection name */
func createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = database.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* GetAll: get all data based on condition*/
func (u *User) GetAll() ([]*User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("users")
	defer cancel()

	//create an empty array to store all fields from collection
	var users []*User = make([]*User, 0)

	//get all user record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to user variable
	for cur.Next(ctx) {
		var user User
		cur.Decode(&user)
		users = append(users, &user)
	}
	//response data to client
	if users == nil {
		return make([]*User, 0), nil
	}
	return users, nil
}

/*GetOne: get one record from a collection  */
func (u *User) GetOne(filter bson.M) (*User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("users")
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

	user := User{}
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
func (u *User) Create(newUser NewUser) (*User, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("users")
	defer cancel()

	//convert newUser to bson.M
	currentTime := time.Now()
	user := User{
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

	return &User{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Email:     newUser.Email,
		Roles:     newUser.Roles}, nil
}

/*UpdateOne: update one record from a collection*/
func (u User) UpdateOne(filter bson.M, update bson.M) (*User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("users")
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
func (u User) DeleteOne(filter bson.M) (*User, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("users")
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

func (u User) Login(input Login) (*User, error) {
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

/* Input Function */
func (u *NewUser) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			user, err := UserInstance.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			if user != nil {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&u.Password, validation.Required.Error("password must not be blanked")),
		validation.Field(&u.ConfirmPassword, validation.Required.Error("confirm password must not be blanked"), validation.In(u.Password).Error("confirm password must be identical with Password")),
		validation.Field(&u.Roles, validation.Required.Error("Role must not be blanked")),
	)
}

func (u *NewUser) HashPassword() error {
	hashPassword, err := utilities.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashPassword
	return nil
}
func (u *Login) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email")),
		validation.Field(&u.Password, validation.Required.Error("password must not be blanked")),
	)
}

//export an instance User model
var UserInstance = User{}
