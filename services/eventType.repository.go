package services

import (
	"context"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var EventTypeRepositoryName = "EventTypeRepositoryName"

type EventTypeRepository struct {
	MongoCN *database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *EventTypeRepository) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* GetAll: get all data based on condition*/
func (u *EventTypeRepository) Find(condition bson.M) ([]*models.EventType, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventTypeName)
	defer cancel()

	//create an empty array to store all fields from collection
	var eventTypes []*models.EventType = make([]*models.EventType, 0)

	//get all record
	cur, err := collection.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var eventType models.EventType
		cur.Decode(&eventType)
		eventTypes = append(eventTypes, &eventType)
	}
	//response data to client
	if eventTypes == nil {
		return make([]*models.EventType, 0), nil
	}
	return eventTypes, nil
}

/*GetOne: get one record from a collection  */
func (u *EventTypeRepository) FindOne(filter bson.M) (*models.EventType, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventTypeName)
	defer cancel()

	eventType := models.EventType{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&eventType); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, helpers.NewErrNotFound("event type id is not found")
		}
		//return err if there is a system error
		return nil, err
	}

	return &eventType, nil
}

/*Create: create a new record to a collection*/
func (u *EventTypeRepository) Create(newEventType model.NewEventType) (*models.EventType, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventTypeName)
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	eventType := models.EventType{
		Name:      newEventType.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		IsDeleted: false,
	}
	newData, err := utilities.InterfaceToBsonM(eventType)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.EventType{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		Name:      newEventType.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		IsDeleted: false,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u EventTypeRepository) UpdateOne(filter bson.M, update bson.M) (*models.EventType, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventTypeName)
	defer cancel()

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, helpers.NewErrNotFound("event type id is not found")
	}

	//query the new update
	eventType, errQuery := u.FindOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return eventType, nil
}

//DeleteOne func is to update one record from a collection
func (u EventTypeRepository) DeleteOne(filter bson.M) (*models.EventType, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventTypeName)
	defer cancel()

	eventType, errGet := u.FindOne(filter)
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
		return nil, helpers.NewErrNotFound("event type id is not found")
	}

	return eventType, nil
}
