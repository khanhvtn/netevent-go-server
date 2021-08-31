package services

import (
	"context"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var FacilityHistoryRepositoryName = "FacilityHistoryRepositoryName"

type FacilityHistoryRepository struct {
	MongoCN *database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *FacilityHistoryRepository) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* FindAll: get all data based on condition*/
func (u *FacilityHistoryRepository) FindAll(condition bson.M) ([]*models.FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionFacilityHistoryName)
	defer cancel()

	//create an empty array to store all fields from collection
	var facilityHistories []*models.FacilityHistory = make([]*models.FacilityHistory, 0)

	//get all record
	cur, err := collection.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var facilityHistory models.FacilityHistory
		cur.Decode(&facilityHistory)
		facilityHistories = append(facilityHistories, &facilityHistory)
	}
	//response data to client
	if facilityHistories == nil {
		return make([]*models.FacilityHistory, 0), nil
	}
	return facilityHistories, nil
}

/*FindOne: get one record from a collection  */
func (u *FacilityHistoryRepository) FindOne(filter bson.M) (*models.FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionFacilityHistoryName)
	defer cancel()

	facilityHistory := models.FacilityHistory{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&facilityHistory); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, helpers.NewErrNotFound("facility history id is not found")
		}
		//return err if there is a system error
		return nil, err
	}

	return &facilityHistory, nil
}

/*Create: create a new record to a collection*/
func (u *FacilityHistoryRepository) Create(newFacilityHistory *models.FacilityHistory) (*models.FacilityHistory, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionFacilityHistoryName)
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	facilityHistory := models.FacilityHistory{
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		Facility:   newFacilityHistory.Facility,
		BorrowDate: newFacilityHistory.BorrowDate,
		ReturnDate: newFacilityHistory.ReturnDate,
		Event:      newFacilityHistory.Event,
	}
	newData, err := utilities.InterfaceToBsonM(facilityHistory)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.FacilityHistory{
		ID:         insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		Facility:   newFacilityHistory.Facility,
		BorrowDate: newFacilityHistory.BorrowDate,
		ReturnDate: newFacilityHistory.ReturnDate,
		Event:      newFacilityHistory.Event,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u FacilityHistoryRepository) UpdateOne(filter bson.M, update bson.M) (*models.FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionFacilityHistoryName)
	defer cancel()

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, helpers.NewErrNotFound("facility history id is not found")
	}

	//query the new update
	facilityHistory, errQuery := u.FindOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return facilityHistory, nil
}

//DeleteOne func is to update one record from a collection
func (u FacilityHistoryRepository) DeleteOne(filter bson.M) (*models.FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionFacilityHistoryName)
	defer cancel()

	facilityHistory, errFind := u.FindOne(filter)
	if errFind != nil {
		return nil, errFind
	}

	//delete user from database
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		//response to client if there is an error.
		return nil, err
	}

	if deleteResult.DeletedCount == 0 {
		return nil, helpers.NewErrNotFound("facility history id is not found")
	}

	return facilityHistory, nil
}
