package database

// import (
// 	"context"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type MongoDBQuery struct {
// 	CollectionName string
// }

// /* createContextAndTargetCol: create and return targeted collection based on collection name */
// func createContextAndTargetCol(collectionName string) (col *mongo.Collection,
// 	ctx context.Context,
// 	cancel context.CancelFunc) {
// 	col = MongoCN.Db.Collection(collectionName)
// 	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 	return
// }

// /* GetAll: get all data based on condition*/
// func (mongoDBQuery *MongoDBQuery) GetAll() (interface{}, *fiber.Error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(mongoDBQuery.CollectionName)
// 	defer cancel()

// 	//create an empty array to store all fields from collection
// 	var data []bson.M

// 	//get all user record
// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		return nil, fiber.NewError(500, err.Error())
// 	}
// 	defer cur.Close(ctx)
// 	//map data to user variable
// 	if err = cur.All(ctx, &data); err != nil {
// 		return nil, fiber.NewError(500, err.Error())
// 	}
// 	//response data to client
// 	if data == nil {
// 		return []bson.M{}, nil
// 	}
// 	return data, nil
// }

// /*GetOne: get one record from a collection  */
// func (mongoDBQuery *MongoDBQuery) GetOne(filter bson.M) (interface{}, *fiber.Error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(mongoDBQuery.CollectionName)
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
// func (mongoDBQuery *MongoDBQuery) Create(newData bson.M) (interface{}, *fiber.Error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(mongoDBQuery.CollectionName)
// 	defer cancel()

// 	//create user in database
// 	insertResult, err := collection.InsertOne(ctx, newData)
// 	if err != nil {
// 		return nil, fiber.NewError(500, err.Error())
// 	}
// 	newData["_id"] = insertResult.InsertedID
// 	return newData, nil
// }

// /*UpdateOne: update one record from a collection*/
// func (mongoDBQuery MongoDBQuery) UpdateOne(filter bson.M, update bson.M) (interface{}, *fiber.Error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(mongoDBQuery.CollectionName)
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
// 	newEventType, errQuery := mongoDBQuery.GetOne(filter)
// 	if errQuery != nil {
// 		return nil, errQuery
// 	}

// 	return newEventType, nil
// }

// //DeleteOne func is to update one record from a collection
// func (mongoDBQuery MongoDBQuery) DeleteOne(filter bson.M) (interface{}, *fiber.Error) {
// 	//get a collection , context, cancel func
// 	collection, ctx, cancel := createContextAndTargetCol(mongoDBQuery.CollectionName)
// 	defer cancel()

// 	result, errGet := mongoDBQuery.GetOne(filter)
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
