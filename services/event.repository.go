package services

import (
	"context"
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var EventRepositoryName = "EventRepositoryName"

type EventRepository struct {
	MongoCN *database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *EventRepository) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* FindAll: get all data based on condition*/
func (u *EventRepository) FindAll(condition *bson.M) ([]*models.Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventName)
	defer cancel()

	//create an empty array to store all fields from collection
	var events []*models.Event = make([]*models.Event, 0)
	//get all record
	var cur *mongo.Cursor
	if condition == nil {
		result, err := collection.Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		cur = result
	} else {
		result, err := collection.Find(ctx, condition)
		if err != nil {
			return nil, err
		}
		cur = result
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var event models.Event
		cur.Decode(&event)
		events = append(events, &event)
	}
	//response data to client
	if events == nil {
		return make([]*models.Event, 0), nil
	}
	return events, nil
}

/*FindOne: get one record from a collection  */
func (u *EventRepository) FindOne(filter bson.M) (*models.Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventName)
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

	event := models.Event{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&event); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &event, nil
}

/*Create: create a new record to a collection*/
func (u *EventRepository) Create(newEvent models.Event) (*models.Event, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventName)
	defer cancel()

	//convert o bson.M
	currentTime := time.Now()
	event := models.Event{
		Tags:                  newEvent.Tags,
		IsApproved:            false,
		Reviewer:              nil,
		IsFinished:            false,
		Tasks:                 newEvent.Tasks,
		FacilityHistories:     newEvent.FacilityHistories,
		Name:                  newEvent.Name,
		Language:              newEvent.Language,
		EventType:             newEvent.EventType,
		Mode:                  newEvent.Mode,
		Location:              newEvent.Location,
		Accommodation:         newEvent.Accommodation,
		RegistrationCloseDate: newEvent.RegistrationCloseDate,
		StartDate:             newEvent.StartDate,
		EndDate:               newEvent.EndDate,
		MaxParticipants:       newEvent.MaxParticipants,
		Description:           newEvent.Description,
		Owner:                 newEvent.Owner,
		Budget:                newEvent.Budget,
		Image:                 newEvent.Image,
		IsDeleted:             false,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
	}
	newData, err := utilities.InterfaceToBsonM(event)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.Event{
		ID:                    insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
		Tags:                  newEvent.Tags,
		IsApproved:            false,
		Reviewer:              nil,
		IsFinished:            false,
		Tasks:                 newEvent.Tasks,
		FacilityHistories:     newEvent.FacilityHistories,
		Name:                  newEvent.Name,
		Language:              newEvent.Language,
		EventType:             newEvent.EventType,
		Mode:                  newEvent.Mode,
		Location:              newEvent.Location,
		Accommodation:         newEvent.Accommodation,
		RegistrationCloseDate: newEvent.RegistrationCloseDate,
		StartDate:             newEvent.StartDate,
		EndDate:               newEvent.EndDate,
		MaxParticipants:       newEvent.MaxParticipants,
		Description:           newEvent.Description,
		Owner:                 newEvent.Owner,
		Budget:                newEvent.Budget,
		Image:                 newEvent.Image,
		IsDeleted:             false}, nil
}

/*UpdateOne: update one record from a collection*/
func (u EventRepository) UpdateOne(filter bson.M, update bson.M) (*models.Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventName)
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
	event, errQuery := u.FindOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return event, nil
}

//DeleteOne func is to update one record from a collection
func (u EventRepository) DeleteOne(filter bson.M) (*models.Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionEventName)
	defer cancel()

	event, errFind := u.FindOne(filter)
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
		return nil, errors.New("id not found")
	}

	return event, nil
}
