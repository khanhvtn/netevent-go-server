package services

// import (
// 	"context"
// 	"time"

// 	"github.com/gofiber/fiber"
// 	"github.com/khanhvtn/netevent-go/database"
// 	"github.com/khanhvtn/netevent-go/models"
// 	"github.com/khanhvtn/netevent-go/utilities"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type UserServiceStruct struct {
// 	collectionName string
// }

// /* createContextAndTargetCol: create and return targeted collection based on collection name */
// func createContextAndTargetCol(collectionName string) (col *mongo.Collection,
// 	ctx context.Context,
// 	cancel context.CancelFunc) {
// 	col = database.MongoCN.Db.Collection(collectionName)
// 	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 	return
// }

// /* GetAll: get all data based on condition*/
// func (uss *UserServiceStruct) GetAll() ([]*models.User, error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(uss.collectionName)
// 	defer cancel()

// 	//create an empty array to store all fields from collection
// 	var users []*models.User = make([]*models.User, 0)

// 	//get all user record
// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		return nil, fiber.NewError(500, err.Error())
// 	}
// 	defer cur.Close(ctx)

// 	//map data to user variable
// 	for cur.Next(ctx) {
// 		var user models.User
// 		cur.Decode(&user)
// 		users = append(users, &user)
// 	}
// 	//response data to client
// 	if users == nil {
// 		return make([]*models.User, 0), nil
// 	}
// 	return users, nil
// }

// /*GetOne: get one record from a collection  */
// func (uss *UserServiceStruct) GetOne(filter bson.M) (interface{}, error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(uss.collectionName)
// 	defer cancel()

// 	//convert id to object id when filter contain _id
// 	if checkID := filter["_id"]; checkID != nil {
// 		if _, ok := checkID.(primitive.ObjectID); !ok {
// 			id, err := primitive.ObjectIDFromHex(checkID.(string))
// 			if err != nil {
// 				return nil, fiber.NewError(500, err.Error())
// 			}
// 			filter["_id"] = id
// 		}
// 	}

// 	result := bson.M{}
// 	//Decode record into result
// 	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
// 		if err != mongo.ErrNoDocuments {
// 			//return err if there is a system error
// 			return nil, fiber.NewError(500, err.Error())
// 		}
// 		//return nil data when id is not existed.
// 		return bson.M{}, fiber.NewError(400, "ID not found")
// 	}

// 	return result, nil
// }

// /*Create: create a new record to a collection*/
// func (uss *UserServiceStruct) Create(newUser models.NewUser) (*models.User, error) {

// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(uss.collectionName)
// 	defer cancel()

// 	//convert newUser to bson.M
// 	user := models.User{
// 		Email:    newUser.Email,
// 		Password: newUser.Password,
// 		Roles:    newUser.Roles,
// 	}
// 	newData, err := utilities.InterfaceToBsonM(user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//create user in database
// 	insertResult, err := collection.InsertOne(ctx, newData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &models.User{ID: insertResult.InsertedID.(primitive.ObjectID),
// 		Email: newUser.Email, Roles: newUser.Roles}, nil
// }

// /*UpdateOne: update one record from a collection*/
// func (uss UserServiceStruct) UpdateOne(filter bson.M, update bson.M) (interface{}, error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(uss.collectionName)
// 	defer cancel()

// 	//convert id to object id when filter contain _id
// 	if checkID := filter["_id"]; checkID != nil {
// 		if _, ok := checkID.(primitive.ObjectID); !ok {
// 			id, err := primitive.ObjectIDFromHex(checkID.(string))
// 			if err != nil {
// 				return nil, fiber.NewError(500, err.Error())
// 			}
// 			filter["_id"] = id
// 		}
// 	}

// 	//update user information
// 	newUpdate := bson.M{"$set": update}
// 	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
// 	if err != nil {
// 		return nil, fiber.NewError(500, err.Error())
// 	}

// 	if updateResult.MatchedCount == 0 {
// 		return bson.M{}, fiber.NewError(400, "ID not found")
// 	}

// 	//query the new update
// 	newEventType, errQuery := uss.GetOne(filter)
// 	if errQuery != nil {
// 		return nil, errQuery
// 	}

// 	return newEventType, nil
// }

// //DeleteOne func is to update one record from a collection
// func (uss UserServiceStruct) DeleteOne(filter bson.M) (interface{}, error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(uss.collectionName)
// 	defer cancel()

// 	result, errGet := uss.GetOne(filter)
// 	if errGet != nil {
// 		return nil, errGet
// 	}

// 	//delete user from database
// 	deleteResult, err := collection.DeleteOne(ctx, filter)
// 	if err != nil {
// 		//response to client if there is an error.
// 		return nil, fiber.NewError(500, err.Error())
// 	}

// 	if deleteResult.DeletedCount == 0 {
// 		return bson.M{}, fiber.NewError(400, "ID not found")
// 	}

// 	return result, nil
// }

// func (uss UserServiceStruct) Login(input models.Login) (*models.User, error) {
// 	return nil, nil
// }

// var UserService = UserServiceStruct{collectionName: "users"}
