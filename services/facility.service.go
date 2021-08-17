package services

import (
	"context"
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var FacilityServiceName = "FacilityServiceName"

type FacilityService struct {
	MongoCN database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *FacilityService) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* GetAll: get all data based on condition*/
func (u *FacilityService) GetAll() ([]*models.Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("facilities")
	defer cancel()

	//create an empty array to store all fields from collection
	var facilities []*models.Facility = make([]*models.Facility, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var facility models.Facility
		cur.Decode(&facility)
		facilities = append(facilities, &facility)
	}
	//response data to client
	if facilities == nil {
		return make([]*models.Facility, 0), nil
	}
	return facilities, nil
}

/*GetOne: get one record from a collection  */
func (u *FacilityService) GetOne(filter bson.M) (*models.Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("facilities")
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

	facility := models.Facility{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&facility); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &facility, nil
}

/*Create: create a new record to a collection*/
func (u *FacilityService) Create(newFacility model.NewFacility) (*models.Facility, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("facilities")
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	facility := model.Facility{
		Name:      newFacility.Name,
		Code:      newFacility.Code,
		Type:      newFacility.Type,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		IsDeleted: false,
		Status:    false,
	}
	newData, err := utilities.InterfaceToBsonM(facility)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.Facility{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      newFacility.Name,
		Code:      newFacility.Code,
		Type:      newFacility.Type,
		IsDeleted: false,
		Status:    false,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u FacilityService) UpdateOne(filter bson.M, update bson.M) (*models.Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("facilities")
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
	facility, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return facility, nil
}

//DeleteOne func is to update one record from a collection
func (u FacilityService) DeleteOne(filter bson.M) (*models.Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("facilities")
	defer cancel()

	facility, errGet := u.GetOne(filter)
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

	return facility, nil
}
